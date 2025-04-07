package server

import (
	"LogGenerator/interfaces"
	"LogGenerator/models"
	"LogGenerator/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)


var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	logGenerationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_generation_total",
			Help: "Total number of log generation tasks",
		},
		[]string{"status"},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(logGenerationCount)
}



// ServerHandler is responsible for handling HTTP requests related to log generation and server health checks.
// It manages incoming requests, including checking the server's health (IsAlive) and initiating log generation tasks (LogHandler).
//
// Fields:
//   - ResponseW: An interface that defines the method for sending HTTP responses.
//     - It is used to send back structured responses to the client.
//   
//   - LogGen: An interface that defines the method for generating logs.
//     - It is used to start log generation tasks and manage concurrent log generation operations.
type ServerHandler struct {
	ResponseW interfaces.ResponseWrite
	LogGen    interfaces.LogGenerator
}

var cancelFunc context.CancelFunc
var mu sync.Mutex // To safely access the cancelFunc in a concurrent environment

// IsAlive handles the "GET /alive" endpoint to check if the server is live.
// It responds with an HTTP status code 200 and a message indicating the server's health status.
//
// Example usage:
//   GET /alive
//   Response: {
//     "status": true,
//     "message": "Server <port> is live",
//     "data": null
//   }
func (s *ServerHandler) IsAlive(w http.ResponseWriter, r *http.Request) {

	timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.Method))
	defer timer.ObserveDuration()

	httpRequests.WithLabelValues(r.Method, "200").Inc()

	s.ResponseW.SendResponse(w, http.StatusOK, true, fmt.Sprintf("Server %v is live", utils.GloablMetaData.Port), nil)
	log.Println("Checking Log Generator Server Call!")
}

// LogHandler handles the "POST /generate" endpoint to initiate log generation.
// It accepts a POST request with a JSON body containing the number of logs to generate and the unit of time (seconds, minutes, or hours).
// After validating the input, it starts a background task to generate the logs and responds with an HTTP status code 200.
// The task will be restarted periodically based on the given duration.
//
// Example usage:
//   POST /generate
//   Request Body: {
//     "num_logs": 1000,
//     "unit": "m"
//   }
//
//   Response: {
//     "status": true,
//     "message": "Task is in progress...",
//     "data": null
//   }
func (s *ServerHandler) LogHandler(w http.ResponseWriter, r *http.Request) {
	response := s.ResponseW
	log.Println("\n Log generation is called!")

	// Start measuring the request duration for Prometheus
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.Method))
	defer timer.ObserveDuration()
	
	// Increment the HTTP request counter for this specific method and status code (e.g., 200)
	httpRequests.WithLabelValues(r.Method, "200").Inc()

	// Default values for rate and unit
	var rate int
	var unitStr string

	var rateModel models.RequestPayload
	if r.Method != http.MethodPost {
		response.SendResponse(w, http.StatusMethodNotAllowed, false, "Only POST method allowed", nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&rateModel)
	if err != nil {
		// Use default values from the config
		rate = int(utils.RateData.NumLogs)
		unitStr = utils.RateData.Unit

		if rate <= 0 || unitStr == "" {
			rate = utils.ConfigData.KEY_RATE
			unitStr = utils.ConfigData.KEY_UNIT
			if rate <= 0 || unitStr == "" {
				response.SendResponse(w, http.StatusBadRequest, false, "Rate and unit are missing", nil)
				return
			}
		}
	} else {
		rate = int(rateModel.NumLogs)
		unitStr = rateModel.Unit
	}

	// Validate unit and set duration
	var duration time.Duration
	switch unitStr {
	case "s":
		duration = 1 * time.Second
	case "m":
		duration = 1 * time.Minute
	case "h":
		duration = 1 * time.Hour
	default:
		response.SendResponse(w, http.StatusBadRequest, false, "Invalid unit. Use s, m, or h for unit variable", nil)
		return
	}

	// Respond immediately with "Task is in progress"
	response.SendResponse(w, http.StatusOK, true, "Task is in progress...", nil)
	log.Println("Response generated to indicate task is in progress")

	// Increment log generation count for this task (started)
	logGenerationCount.WithLabelValues("started").Inc()

	// Cancel the previous task if it exists
	mu.Lock()
	if cancelFunc != nil {
		cancelFunc() // Cancel the previous task
		log.Println("Previous task canceled.")
	}
	mu.Unlock()

	// Start the background task after responding
	go s.startLogGenerationTask(rate, unitStr, duration)
}

// startLogGenerationTask starts the log generation task in the background.
// It runs periodically based on the specified duration (rate and unit) and can be canceled or restarted as needed.
//
// Fields:
//   - rate: The number of logs to generate during each period.
//   - unitStr: The unit of time for the task's duration (either "s", "m", or "h").
//   - duration: The duration between each log generation task. It is calculated based on the unit provided.
//
// It starts a background task to generate logs and cancels the previous task if it's still running.
func (s *ServerHandler) startLogGenerationTask(rate int, unitStr string, duration time.Duration) {
	// Create a new context for the current task
	cntx, cancel := context.WithCancel(context.Background())
	mu.Lock()
	cancelFunc = cancel
	mu.Unlock()

	var wg sync.WaitGroup
	ticker := time.NewTicker(duration)

	// Start the log generation task in a goroutine
	go s.LogGen.GenerateLogsConcurrently(cntx, rate, duration, &wg)

	for {
		select {
		case <-ticker.C:
			// Cancel the previous task and start a new one
			mu.Lock()
			if cancelFunc != nil {
				cancelFunc()
			}
			cntx, cancel = context.WithCancel(context.Background())
			cancelFunc = cancel
			mu.Unlock()

			wg.Add(1)
			log.Println("-------------------------------------------------------")
			go s.LogGen.GenerateLogsConcurrently(cntx, rate, duration, &wg)

		case <-cntx.Done():
			// Task was externally stopped (e.g., by cancelFunc)
			log.Println("Stopped externally")
			return
		}
	}

	// Optionally, you can wait for the tasks to complete if needed
	// wg.Wait()
}

// Expose /metrics endpoint for Prometheus
func (s *ServerHandler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the /metrics endpoint
	promhttp.Handler().ServeHTTP(w, r)
}