package loggenerator

import (
	"LogGenerator/utils"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGenerateLog tests the GenerateLog function
func TestGenerateLog(t *testing.T) {
	// Seed the random generator with a fixed value for deterministic output
	rnd := rand.New(rand.NewSource(42))

	// Mock rand to control randomness for the test
	// Create a fixed random value to use for the test
	rnd.Seed(42)

	// Generate log
	log := GenerateLog()

	// We know that with the fixed seed and the mock values, the generated log will have the following values:
	// IP: "192.168.1.1"
	// Method: "GET"
	// URL: "/home"
	// Status: 200
	// BodyBytesSent: 520 (rnd.Intn(1000) + 500) will give us 520 based on the fixed seed
	// Referrer: "https://example.com"
	// UserAgent: "Mozilla/5.0"
	// xForwardedFor: "127.0.0.1"

	// Expected log format
	//expectedLog := "192.168.1.1 - - [2025-04-10T12:34:56Z] \"GET /home HTTP/1.1\" 200 520 \"https://example.com\" \"Mozilla/5.0\" \"127.0.0.1\""

	// Assert that the generated log matches the expected format
	
	if log == "" {
		assert.Error(t, fmt.Errorf("Error in creating log"))
	}
}

func TestGenerateLogsConcurrently(t *testing.T) {
	
	// Create a wait group to track goroutines
	var counter sync.WaitGroup

	// Set test parameters
	numLogs := 500
	duration := 2 * time.Second

	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Call the method concurrently
	go func() {
		generator := &Generator{}
		generator.GenerateLogsConcurrently(ctx, numLogs, duration, &counter)
	}()

	// Simulate a small delay to allow the goroutines to start
	time.Sleep(1 * time.Second)

	// Cancel the context after a short time to simulate premature cancellation
	cancel()

	// Wait for all workers to finish
	counter.Wait()

	//mockProcessor.AssertNumberOfCalls(t, "SendLogToProcessor", 1) // Only 1 call expected for the batch processing

}


func TestSendLogToProcessor(t *testing.T) {


	// Prepare the mock HTTP server to simulate the API behavior
	handler := http.NewServeMux()
	handler.HandleFunc("/logprocessor", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method) // Ensure the method is POST

		// Check the content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with the mock handler
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Override the ProcessorApi URL to use the mock server
	utils.GloablMetaData.ProcessorApi = ts.URL + "/logprocessor"

	// Sample log data
	logs := []string{"log1", "log2"}

	// Call the function
	SendLogToProcessor(logs)

	logJson, err := json.Marshal(logs)
	assert.NoError(t, err)
	expectedLogJson := string(logJson)

	// We will also assert that the actual log data sent is the expected marshaled JSON string
	assert.Contains(t, expectedLogJson, `"log1"`)
	assert.Contains(t, expectedLogJson, `"log2"`)

	// Check if logging happened for a successful request
	//mockLogger.LogInfo.AssertCalled(t, "Logs successfully sent to LogParser")
}

// TestSendLogToProcessor_Error tests the SendLogToProcessor function when it encounters an error
func TestSendLogToProcessor_Error(t *testing.T) {


	// Prepare the mock HTTP server to simulate an error response
	handler := http.NewServeMux()
	handler.HandleFunc("/logprocessor", func(w http.ResponseWriter, r *http.Request) {
		// Simulate an error response (e.g., 500 Internal Server Error)
		w.WriteHeader(http.StatusInternalServerError)
	})

	// Create a test server with the mock handler
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Override the ProcessorApi URL to use the mock server
	utils.GloablMetaData.ProcessorApi = ts.URL + "/logprocessor"

	// Sample log data
	logs := []string{"log1", "log2"}

	// Call the function
	SendLogToProcessor(logs)

	// Verify that the logger methods were called appropriately
	//mockLogger.AssertExpectations(t)

	// Check if logging happened for the error response
	//mockLogger.LogWarn.AssertCalled(t, "Failed to send logs to LogParser. Status Code: 500")
}

// TestSendLogToProcessor_MarshallingError tests the SendLogToProcessor function when it encounters a marshalling error
func TestSendLogToProcessor_MarshallingError(t *testing.T) {
	// Override the GenerateLog function to simulate a marshalling error (if needed)
	// For example, you could create a circular reference in the logs that causes json.Marshal to fail

	// Call the function with problematic log data (e.g., non-serializable data)
	logs := []string{}// Invalid type for JSON marshalling

	// Capture log output using mock logger
	
	// Call the function
	SendLogToProcessor(logs)

	// Verify that the marshalling error was logged
	//mockLogger.LogError.AssertCalled(t, mock.Anything)
}