package handlers

import (
	"LogParser/connection"
	"LogParser/models"
	"LogParser/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Prometheus metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)
	logAddsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "log_adds_total",
			Help: "Total number of logs added",
		},
	)
	httpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(logAddsTotal)
	prometheus.MustRegister(httpRequestDurationSeconds)
}

// IsAlive checks if the server is running and responds with an HTTP 200 OK status.
func IsAlive(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(httpRequestDurationSeconds.WithLabelValues("GET"))
	defer timer.ObserveDuration()

	models.SendResponse(w, http.StatusOK, true, fmt.Sprintf("Server %v is live", utils.ConfigData.PORT),nil)
	log.Println("checking the server call!")

	httpRequestsTotal.WithLabelValues("GET", "200").Inc()
}

// HandleType handles HTTP requests based on the method type (POST, GET, DELETE).
func HandleType(w http.ResponseWriter, r *http.Request){
	timer := prometheus.NewTimer(httpRequestDurationSeconds.WithLabelValues(r.Method))
	defer timer.ObserveDuration()

	switch r.Method{
	case http.MethodPost:
		AddLogsHandler(w,r)
	case http.MethodGet:
		GetLogsHandler(w,r)
	case http.MethodDelete:
		DeleteLogsHandler(w,r)
	default:
		log.Println("Method not allowed!")
		models.SendResponse(w, http.StatusMethodNotAllowed, false, "Only GET, POST, DELETE methods are allowed to execute the task", nil)
		httpRequestsTotal.WithLabelValues(r.Method, "405").Inc()
		//GetLogsHandler(w,r)
	}
}

// GetLogsCountHandler returns the count of logs based on the applied filters.
func GetLogsCountHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get logs count hit!")
	timer := prometheus.NewTimer(httpRequestDurationSeconds.WithLabelValues("GET"))
	defer timer.ObserveDuration()

	isAlive, db := connection.PingDB()
	if !isAlive {
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to connect to Database!"), nil)
		httpRequestsTotal.WithLabelValues("GET", "500").Inc()
		return
	}

	dateFilter, _ := utils.GetDateFilters(r)
	query, args := utils.GenerateFilteredCountQuery(utils.GenerateFiltersMap(r), utils.GetPaginationParams(r), dateFilter)

	var count int
	err := db.QueryRow(query, args...).Scan(&count)

	if err != nil {
		// If there's an error querying the database
		log.Printf("Failed to query database: %v", err)
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to query database: %v", err), nil)
		httpRequestsTotal.WithLabelValues("GET", "500").Inc()
		return
	}

	if count <= 0 {
		models.SendResponse(w, http.StatusOK, true, "No logs found", nil)
		httpRequestsTotal.WithLabelValues("GET", "200").Inc()
	} else {
		models.SendResponse(w, http.StatusOK, true, "Logs count fetched successfully", map[string]int{"count": count})
		httpRequestsTotal.WithLabelValues("GET", "200").Inc()
	}
}

// GetLogsHandler fetches logs based on filters and pagination, and returns them in the response.
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(httpRequestDurationSeconds.WithLabelValues("GET"))
	defer timer.ObserveDuration()
	isAlive, db := connection.PingDB()
	if !isAlive {
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to connect to Database!"), nil)
		httpRequestsTotal.WithLabelValues("GET", "500").Inc()
		return
	}
	dateFilter, _ := utils.GetDateFilters(r)
	query, args := utils.GenerateFilteredGetQuery(utils.GenerateFiltersMap(r), utils.GetPaginationParams(r), dateFilter)

	//query, args := utils.GenerateFilteredGetQuery(utils.GenerateFiltersMap(r))
	rows, err := db.Query(query, args...)

	if err != nil {
		log.Printf("Failed to query database: %v", err)
		models.SendResponse(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("Failed to query database : %v", err), nil)
		httpRequestsTotal.WithLabelValues("GET", "500").Inc()
		return

	}

	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		if err := rows.Scan(&log.RemoteAddr, &log.RemoteUser, &log.TimeLocal, &log.Request, &log.Status, &log.BodyBytesSent, &log.HttpReferer, &log.HttpUserAgent, &log.HttpXForwardedFor); err != nil {
			fmt.Printf("Failed to scan log: %v", err)
			models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to scan log : %v", err), nil)
			httpRequestsTotal.WithLabelValues("GET", "500").Inc()
			return
		}
		logs = append(logs, log)
	}

	if len(logs) == 0 {
		models.SendResponse(w, http.StatusOK, true, "No Logs found", logs)
		httpRequestsTotal.WithLabelValues("GET", "200").Inc()
	} else {
		models.SendResponse(w, http.StatusOK, true, fmt.Sprintf("%d Logs fetched successfully",len(logs)), logs)
		httpRequestsTotal.WithLabelValues("GET", "200").Inc()
	}
}

// DeleteLogsHandler deletes logs from the database based on the filters provided in the request.
func DeleteLogsHandler(w http.ResponseWriter, r *http.Request) {
	isAlive, db := connection.PingDB()
	if !isAlive {
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to connect to Database!", nil)
		return
	}

	query, args := utils.GenerateDeleteQuery(utils.GenerateFiltersMap(r))

	result, err := db.Exec(query, args...)
	if err != nil {
		// Log error and send response if the query fails
		log.Printf("Failed to execute delete query: %v", err)
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to execute delete query: %v", err), nil)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get affected rows: %v", err)
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to get affected rows: %v", err), nil)
		return
	}

	if rowsAffected > 0 {
		models.SendResponse(w, http.StatusOK, true, fmt.Sprintf("%d logs deleted successfully.", rowsAffected), nil)
	} else {
		models.SendResponse(w, http.StatusOK, true, "No logs found matching the provided filters.", nil)
	}
}

// InsertOneLog inserts a single log entry into the database.
func InsertOneLog(logs models.Log) error {
	isAlive, db := connection.PingDB()
	if !isAlive {
		return fmt.Errorf("Database is down!")
	}
	_, err := db.Exec(`INSERT INTO logs (remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, logs.RemoteAddr, logs.RemoteUser, logs.TimeLocal, logs.Request, logs.Status, logs.BodyBytesSent, logs.HttpReferer, logs.HttpUserAgent, logs.HttpXForwardedFor)

	if err != nil {
		log.Printf("Error inserting log: %v", err) // More detailed error logging
		return err
	}
	return nil
}

// AddLogsHandler processes the incoming POST request and inserts logs into the database.
func AddLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Add hit!")
	timer := prometheus.NewTimer(httpRequestDurationSeconds.WithLabelValues("POST"))
	defer timer.ObserveDuration()

	if r.Method != http.MethodPost {
		models.SendResponse(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("%d Invalid request method", http.StatusMethodNotAllowed), nil)
		httpRequestsTotal.WithLabelValues("POST", "405").Inc()
		return
	}

	var logstr []string
	err := json.NewDecoder(r.Body).Decode(&logstr)
	if err != nil {
		http.Error(w, "Failed to decode log data", http.StatusBadRequest)
		log.Printf("Error decoding log data: %v", err)
		httpRequestsTotal.WithLabelValues("POST", "400").Inc()
		return
	}

	isAlive, db := connection.PingDB()
	if !isAlive {
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to connect to Database!", nil)
		httpRequestsTotal.WithLabelValues("POST", "500").Inc()
		return
	}

	count := len(logstr)
	log.Println("Received : ",count)
	
	logsChan := make(chan string, len(logstr))
	resultsChan := make(chan models.Log, len(logstr))

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU() 
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processLogWorker(logsChan, resultsChan, &wg)
	}

	for _, logStr := range logstr {
		logsChan <- logStr
	}
	close(logsChan)

	go func() {
		wg.Wait()
		close(resultsChan) 
	}()

	var logEntries []models.Log
	for logEntry := range resultsChan {
		logEntries = append(logEntries, logEntry)
	}

	query, values := utils.GenerateAddQuery(logEntries)
	result, err1 := db.Exec(query, values...)
	if err1 != nil {
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to insert logs: %v", err1), nil)
		log.Printf("Failed to insert logs: %v", err1)
		httpRequestsTotal.WithLabelValues("POST", "500").Inc()
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to retrieve affected rows: %v", err), nil)
		log.Printf("Error retrieving affected rows: %v", err)
		httpRequestsTotal.WithLabelValues("POST", "500").Inc()
		return
	}

	models.SendResponse(w, http.StatusOK, true, fmt.Sprintf("Logs stored successfully, %d rows inserted.", rowsAffected), nil)
	httpRequestsTotal.WithLabelValues("POST", "200").Inc()
}

// MetricsHandler exposes the /metrics endpoint for Prometheus scraping.
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// processLogWorker processes logs concurrently, transforming log strings into log entries.
func processLogWorker(logs <-chan string, results chan<- models.Log, wg *sync.WaitGroup) {
	defer wg.Done()
	for logStr := range logs {
		logEntry := parseLog(logStr)
		results <- logEntry
	}
}

// parseLog converts a log string into a structured Log entry using regular expressions.
func parseLog(logStr string) models.Log {
	// log format example:
	// 192.168.1.2 - - [17/Mar/2025:13:30:20 +0530] "GET /home HTTP/1.1" 500 1180 "https://www.bing.com" "Mozilla/5.0..."

	re := regexp.MustCompile(`(?P<RemoteAddr>[\d\.]+) - (?P<RemoteUser>[-\w]+) \[([^\]]+)\] "(?P<Request>[^"]+)" (?P<Status>\d+) (?P<BodyBytesSent>\d+) "(?P<HttpReferer>[^"]*)" "(?P<HttpUserAgent>[^"]*)" "(?P<HttpXForwardedFor>[^"]*)"`)
	matches := re.FindStringSubmatch(logStr)

	if len(matches) > 0 {
		return models.Log{
			RemoteAddr:    matches[1],
			RemoteUser:    matches[2],
			TimeLocal:     matches[3],
			Request:       matches[4],
			Status:        atoi(matches[5]),
			BodyBytesSent: atoi(matches[6]),
			HttpReferer:   matches[7],
			HttpUserAgent: matches[8],
			HttpXForwardedFor: matches[9],
		}
	}
	return models.Log{}
}

// atoi safely converts a string to an integer, returning 0 on error.
func atoi(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}