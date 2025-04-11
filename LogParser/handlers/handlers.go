package handlers

import (
	"LogParser/connection"
	"LogParser/logger"
	"LogParser/models"
	"LogParser/utils"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// IsAlive checks if the server is running and responds with an HTTP 200 OK status.
func IsAlive(w http.ResponseWriter, r *http.Request) {
	models.SendResponse(w, http.StatusOK, true, fmt.Sprintf("Server %v is live", utils.ConfigData.PORT),nil)
	logger.LogDebug("checking the server call!")
}

// HandleType handles HTTP requests based on the method type (POST, GET, DELETE).
func HandleType(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodPost:
		AddLogsHandler(w,r)
	case http.MethodGet:
		GetLogsHandler(w,r)
	case http.MethodDelete:
		DeleteLogsHandler(w,r)
	default:
		logger.LogWarn("Method not allowed!")
		models.SendResponse(w, http.StatusMethodNotAllowed, false, "Only GET, POST, DELETE methods are allowed to execute the task", nil)
		//GetLogsHandler(w,r)
	}
}

// GetLogsCountHandler returns the count of logs based on the applied filters.
func GetLogsCountHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogDebug("Get logs count hit!")

	isAlive, db := connection.PingDB()
	if !isAlive {
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to connect to Database!"), nil)
		return
	}

	var totalLogs int
	err := db.QueryRow(utils.QUERY_COUNT_ALL).Scan(&totalLogs)
	if err != nil {
		logger.LogWarn(fmt.Sprintf("Error fetching total log count: %v", err))
	}

	//dateFilter, _ := utils.GetDateFilters(r)
	query, args := utils.GenerateFilteredCountQuery(utils.GenerateFiltersMap(r))//, utils.GetPaginationParams(r), dateFilter

	var count int
	err1 := db.QueryRow(query, args...).Scan(&count)
	if err1 != nil {
		logger.LogWarn(fmt.Sprintf("Failed to query database: %v", err1))
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to query database: %v", err1), nil)
		return
	}

	if count <= 0 {
		models.SendResponse(w, http.StatusOK, true, "No logs found", nil)
	} else {
		data := map[string]int{
			"total": totalLogs,
			"fetch": count,
		}
		models.SendResponse(w, http.StatusOK, true, "Logs Found Success", data)
	}
}

// GetLogsHandler fetches logs based on filters and pagination, and returns them in the response.
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get logs api hit!")
	isAlive, db := connection.PingDB()
	if !isAlive {
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to connect to Database!"), nil)
		return
	}

	var totalLogs int
	err := db.QueryRow(utils.QUERY_COUNT_ALL).Scan(&totalLogs)
	if err != nil {
		logger.LogWarn(fmt.Sprintf("Error fetching total log count: %v", err))
	}

	dateFilter, errs := utils.GetDateFilters(r)
	if errs != nil {
		logger.LogWarn(fmt.Sprintf("Error in parsing filetered dates:%v", errs))
	}
	
	paginationFilter := utils.GetPaginationParams(r)
	query, args := utils.GenerateFilteredGetQuery(utils.GenerateFiltersMap(r), paginationFilter, dateFilter)

	//query, args := utils.GenerateFilteredGetQuery(utils.GenerateFiltersMap(r))
	rows, err := db.Query(query, args...)

	if err != nil {
		logger.LogWarn(fmt.Sprintf("Failed to query database: %v", err))
		models.SendResponse(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("Failed to query database : %v", err), nil)
		return

	}
	defer rows.Close()

	var logs []models.Log

	var firstLogTime, lastLogTime time.Time
	for rows.Next() {
		var log models.Log
		if err := rows.Scan(&log.RemoteAddr, &log.RemoteUser, &log.TimeLocal, &log.Request, &log.Status, &log.BodyBytesSent, &log.HttpReferer, &log.HttpUserAgent, &log.HttpXForwardedFor); err != nil {
			fmt.Printf("Failed to scan log: %v", err)
			models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to scan log : %v", err), nil)
			return
		}
		logs = append(logs, log)

		// Set the first and last log times for pagination.
		if firstLogTime.IsZero() {
			firstLogTime = log.TimeLocal
		}
		lastLogTime = log.TimeLocal
	}

	var nextCursor, prevCursor *string
	if len(logs) > 0 {
		// The nextCursor will be the last log's time if it's not the last page
		if len(logs) == paginationFilter.Limit {
			nextCursor = FormatTime(&lastLogTime)  // Only set nextCursor if we have more logs to fetch
		}

		// The prevCursor will be the first log's time if it's not the first page
		if paginationFilter.Cursor != nil && len(logs) > 0 {
			prevCursor = FormatTime(&firstLogTime)
		}
	}

	responseData := map[string]interface{}{
		"count": map[string]interface{}{
			"total": totalLogs,    // Total logs in the database
			"fetch": len(logs),    // Logs fetched in this request
		},
		"logs": logs, // Logs fetched from the database
		"paging": map[string]interface{}{
			"next_cursor": nextCursor,  // Cursor for next page
			"prev_cursor": prevCursor,  // Cursor for previous page
			"limit": paginationFilter.Limit, // The limit used for pagination
		},
	}

	if len(logs) == 0 {
		
		models.SendResponse(w, http.StatusOK, true, "No Logs found", responseData)
	} else {
		models.SendResponse(w, http.StatusOK, true, "Fetched logs successfully", responseData)
	}
}

func FormatTime(t *time.Time) *string {
    if t == nil {
        return nil
    }
    formattedTime := t.Format(time.RFC3339)
    return &formattedTime
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
		logger.LogWarn(fmt.Sprintf("Failed to execute delete query: %v", err))
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to execute delete query: %v", err), nil)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogWarn(fmt.Sprintf("Failed to get affected rows: %v", err))
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
		logger.LogError(fmt.Sprintf("Error inserting log: %v", err)) // More detailed error logging
		return err
	}
	return nil
}

// AddLogsHandler processes the incoming POST request and inserts logs into the database.
func AddLogsHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogDebug("Add hit!")

	if r.Method != http.MethodPost {
		models.SendResponse(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("%d Invalid request method", http.StatusMethodNotAllowed), nil)
		return
	}

	var logstr []string
	err := json.NewDecoder(r.Body).Decode(&logstr)
	if err != nil {
		http.Error(w, "Failed to decode log data", http.StatusBadRequest)
		logger.LogError(fmt.Sprintf("Error decoding log data: %v", err))
		return
	}

	isAlive, db := connection.PingDB()
	if !isAlive {
		models.SendResponse(w, http.StatusInternalServerError, false, "Failed to connect to Database!", nil)
		return
	}

	count := len(logstr)
	logger.LogDebug(fmt.Sprintf("Received : %v",count))
	
	logsChan := make(chan string, len(logstr))
	resultsChan := make(chan models.Log, len(logstr))

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU() 
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go ProcessLogWorker(logsChan, resultsChan, &wg)
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
		logger.LogWarn(fmt.Sprintf("Failed to insert logs: %v", err1))
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		models.SendResponse(w, http.StatusInternalServerError, false, fmt.Sprintf("Failed to retrieve affected rows: %v", err), nil)
		logger.LogError(fmt.Sprintf("Error retrieving affected rows: %v", err))
		return
	}

	models.SendResponse(w, http.StatusOK, true, fmt.Sprintf("Logs stored successfully, %d rows inserted.", rowsAffected), nil)
}

// processLogWorker processes logs concurrently, transforming log strings into log entries.
func ProcessLogWorker(logs <-chan string, results chan<- models.Log, wg *sync.WaitGroup) {
	defer wg.Done()
	for logStr := range logs {
		logEntry := ParseLog(logStr)
		results <- logEntry
	}
}

func ParseLog(logStr string) models.Log {
	// Define a regular expression to capture the log fields
	re := regexp.MustCompile(`^([\d\.]+) - (\S+) \[([^\]]+)\] "(.*?)" (\d{3}) (\d+) "(.*?)" "(.*?)" "(.*?)"$`)
	matches := re.FindStringSubmatch(logStr)

	if len(matches) > 0 {
		// Parse the time field into a time.Time object
		logTime, err := time.Parse(time.RFC3339, matches[3])
		if err != nil {
			logTime = time.Time{} // Default to zero time if parsing fails
		}

		// Return a structured Log model
		return models.Log{
			RemoteAddr:       matches[1],
			RemoteUser:       matches[2],
			TimeLocal:        logTime, // Store as time.Time
			Request:          matches[4],
			Status:           Atoi(matches[5]),
			BodyBytesSent:    Atoi(matches[6]),
			HttpReferer:      matches[7],
			HttpUserAgent:    matches[8],
			HttpXForwardedFor: matches[9],
		}
	}

	// Return empty log if the format doesn't match
	return models.Log{}
}

/*
func parseLog2(logEntry string) (models.Log, error) {
	// Define a regex pattern to match the log structure
	// The pattern captures each part of the log entry.
	// Example log: "10.0.0.1 - - [2025-04-08T06:57:31Z] \"GET /login HTTP/1.1\" 301 1043 \"https://www.bing.com\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/18.18362\" \"212.32.188.247\""
	re := regexp.MustCompile(`^(\S+) (\S+) (\S+) \[([^\]]+)\] "([A-Z]+) (.*?) HTTP/\S+" (\d{3}) (\d+) "(.*?)" "(.*?)" "(.*?)"$`)

	// Match the log entry against the regex pattern
	matches := re.FindStringSubmatch(logEntry)
	if len(matches) < 11 {
		return models.Log{}, fmt.Errorf("invalid log format")
	}

	// Extract fields from the regex matches
	remoteAddr := matches[1]
	remoteUser := matches[2]
	timeStr := matches[4]
	request := matches[5] + " " + matches[6] + " HTTP/1.1"
	status := matches[7]
	bodyBytesSent := matches[8]
	httpReferer := matches[9]
	httpUserAgent := matches[10]
	httpXForwardedFor := matches[11]

	// Parse time
	logTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return models.Log{}, fmt.Errorf("invalid time format: %v", err)
	}

	// Convert status and bodyBytesSent to integers
	statusCode, err := strconv.Atoi(status)
	if err != nil {
		return models.Log{}, fmt.Errorf("invalid status code: %v", err)
	}

	bodyBytes, err := strconv.Atoi(bodyBytesSent)
	if err != nil {
		return models.Log{}, fmt.Errorf("invalid body bytes sent: %v", err)
	}

	// Create Log struct and return it
	log := models.Log{
		RemoteAddr:       remoteAddr,
		RemoteUser:       remoteUser,
		TimeLocal:        logTime,
		Request:          request,
		Status:           statusCode,
		BodyBytesSent:    bodyBytes,
		HttpReferer:      httpReferer,
		HttpUserAgent:    httpUserAgent,
		HttpXForwardedFor: httpXForwardedFor,
	}

	return log, nil
}

// parseLog converts a log string into a structured Log entry using regular expressions.
func parseLog1(logStr string) models.Log {
	// log format example:
	// 192.168.1.2 - - [17/Mar/2025:13:30:20 +0530] "GET /home HTTP/1.1" 500 1180 "https://www.bing.com" "Mozilla/5.0..."

	re := regexp.MustCompile(`(?P<RemoteAddr>[\d\.]+) - (?P<RemoteUser>[-\w]+) \[([^\]]+)\] "(?P<Request>[^"]+)" (?P<Status>\d+) (?P<BodyBytesSent>\d+) "(?P<HttpReferer>[^"]*)" "(?P<HttpUserAgent>[^"]*)" "(?P<HttpXForwardedFor>[^"]*)"`)
	matches := re.FindStringSubmatch(logStr)

	if len(matches) > 0 {
		return models.Log{
			RemoteAddr:    matches[1],
			RemoteUser:    matches[2],
			//TimeLocal:     matches[3],
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
	*/

// atoi safely converts a string to an integer, returning 0 on error.
func Atoi(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}