package loggenerator

import (
	"LogGenerator/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)


// SendLogToProcessor sends a batch of logs to an external log processor via an HTTP POST request.
// The logs are sent in JSON format to the log processor API endpoint specified in the configuration.
//
// Parameters:
//   - logs: A slice of strings containing the log entries to be sent to the processor.
//     These logs are marshaled into JSON format before being sent in the request body.
//
// The function does the following:
//   1. Marshals the logs into a JSON format.
//   2. Creates a new HTTP client with a timeout of 10 seconds.
//   3. Sends an HTTP POST request to the log processor API, including the marshaled logs in the body.
//   4. Handles potential errors, logs the results, and prints success/failure messages based on the HTTP response.
//
// If the request is successful (HTTP status 200 OK), it logs a success message.
// If there's any error (either in marshalling or the HTTP request), it logs the error details.
//
// Example usage:
//   logs := []string{"log1", "log2", "log3"}
//   SendLogToProcessor(logs)
func SendLogToProcessor(logs []string) {
	log.Println("Send log is called!")
	logJson, err := json.Marshal(logs)
	if err != nil {
		log.Printf("Error marshalling log data: %v", err)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second, 
	}

	resp, err := client.Post(utils.GloablMetaData.ProcessorApi, "application/json", bytes.NewBuffer(logJson))
	if err != nil {
		//log.Println("Log processor error : ",err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Logs successfully sent to LogParser")
	} else {
		log.Printf("Failed to send logs to LogParser. Status Code: %d", resp.StatusCode)
	}
}