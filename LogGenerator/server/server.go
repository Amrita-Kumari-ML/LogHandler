package server

import (
	"LogGenerator/interfaces"
	"LogGenerator/logger"
	"LogGenerator/models"
	"LogGenerator/utils"
	"context"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"sync"
	"time"
)

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
var mu sync.Mutex 

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
	s.ResponseW.SendResponse(w, http.StatusOK, true, fmt.Sprintf("Server %v is live", utils.GloablMetaData.Port), nil)
	logger.LogDebug("Checking Log Generator Server Call!")
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
	logger.LogDebug("\n Log generation is called!")

	var rate int
	var unitStr string

	var rateModel models.RequestPayload
	if r.Method != http.MethodPost {
		response.SendResponse(w, http.StatusMethodNotAllowed, false, "Only POST method allowed", nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&rateModel)
	if err != nil {
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

	statusChan := make(chan string, 1) // Buffered so it doesn't block
	mu.Lock()
	if cancelFunc != nil {
		cancelFunc() 
		logger.LogWarn("Previous task canceled.")
	}
	mu.Unlock()

	go s.startLogGenerationTask(rate, unitStr, duration, statusChan)

	select {
	case statusMsg := <-statusChan:
		response.SendResponse(w, http.StatusOK, true, statusMsg, nil)
		logger.LogInfo("Response generated to indicate task is in progress")
	case <-time.After(3 * time.Second):
		response.SendResponse(w, http.StatusRequestTimeout, false, "No status received in time", nil)
		logger.LogWarn("No status received in time")
	}
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
func (s *ServerHandler) startLogGenerationTask(rate int, unitStr string, duration time.Duration, statusChan chan<- string) {
	cntx, cancel := context.WithCancel(context.Background())
	mu.Lock()
	cancelFunc = cancel
	mu.Unlock()

	var wg sync.WaitGroup
	ticker := time.NewTicker(duration)
	if rate <= 0 {
		msg := fmt.Sprintf("numLogs is zero or negative, skipping the generate")
		logger.LogError(msg)
		select {
		case statusChan <- msg:
		default:
		}
		cntx.Done()
		return
	}
	go s.LogGen.GenerateLogsConcurrently(cntx, rate, duration, &wg, statusChan)

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			if cancelFunc != nil {
				cancelFunc()
			}
			cntx, cancel = context.WithCancel(context.Background())
			cancelFunc = cancel
			mu.Unlock()

			wg.Add(1)
			go s.LogGen.GenerateLogsConcurrently(cntx, rate, duration, &wg, statusChan)

		case <-cntx.Done():
			logger.LogWarn("Stopped externally")
			return
		}
	}

	// Optionally, you can wait for the tasks to complete if needed
	// wg.Wait()
}
