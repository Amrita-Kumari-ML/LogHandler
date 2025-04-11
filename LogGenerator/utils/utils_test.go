// utils/utils_test.go
package utils

import (
	"LogGenerator/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	// Test the configuration keys constants
	assert.Equal(t, KEY_HOST, "GENERATOR_HOST", "KEY_HOST should be 'GENERATOR_HOST'")
	assert.Equal(t, KEY_PORT, "GENERATOR_PORT", "KEY_PORT should be 'GENERATOR_PORT'")
	assert.Equal(t, KEY_ALIVE_URL, "GENERATOR_ALIVE_URL", "KEY_ALIVE_URL should be 'GENERATOR_ALIVE_URL'")
	assert.Equal(t, KEY_START_URL, "GENERATOR_START_URL", "KEY_START_URL should be 'GENERATOR_START_URL'")
	assert.Equal(t, KEY_PARSER_API, "PARSER_API", "KEY_PARSER_API should be 'PARSER_API'")
	assert.Equal(t, KEY_RATE, "GENERATOR_RATE", "KEY_RATE should be 'GENERATOR_RATE'")
	assert.Equal(t, KEY_UNIT, "GENERATOR_UNIT", "KEY_UNIT should be 'GENERATOR_UNIT'")

	// Test the default values
	assert.Equal(t, GENERATOR_HOST, "loggenerate", "GENERATOR_HOST should be 'loggenerate'")
	assert.Equal(t, GENERATOR_PORT, ":8080", "GENERATOR_PORT should be ':8080'")
	assert.Equal(t, GENERATOR_ALIVE_URL, "/", "GENERATOR_ALIVE_URL should be '/'")
	assert.Equal(t, GENERATOR_START_URL, "/logs", "GENERATOR_START_URL should be '/logs'")
	assert.Equal(t, PARSER_API, "http://localhost:8083/logs", "PARSER_API should be 'http://localhost:8083/logs'")
	assert.Equal(t, GENERATOR_RATE, 10, "GENERATOR_RATE should be 10")
	assert.Equal(t, GENERATOR_UNIT, "s", "GENERATOR_UNIT should be 's'")
	assert.Equal(t, FILE_NAME, "config.yaml", "FILE_NAME should be 'config.yaml'")
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestIps(t *testing.T) {
	// Verify the contents of the Ips slice
	expectedIps := []string{
		"192.168.1.1", 
		"192.168.1.2", 
		"10.0.0.1",
	}
	assert.Equal(t, Ips, expectedIps, "Ips slice does not match expected values")
}

func TestMethods(t *testing.T) {
	// Verify the contents of the Methods slice
	expectedMethods := []string{
		"GET", 
		"POST", 
		"PUT", 
		"DELETE",
	}
	assert.Equal(t, Methods, expectedMethods, "Methods slice does not match expected values")
}

func TestUrls(t *testing.T) {
	// Verify the contents of the Urls slice
	expectedUrls := []string{
		"/home", 
		"/login", 
		"/profile", 
		"/dashboard",
	}
	assert.Equal(t, Urls, expectedUrls, "Urls slice does not match expected values")
}

func TestStatuses(t *testing.T) {
	// Verify the contents of the Statuses slice
	expectedStatuses := []int{
		200, // OK
		404, // Not Found
		500, // Internal Server Error
		301, // Moved Permanently
	}
	assert.Equal(t, Statuses, expectedStatuses, "Statuses slice does not match expected values")
}

func TestUserAgents(t *testing.T) {
	// Verify the contents of the UserAgents slice
	expectedUserAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/18.18362",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
	}
	assert.Equal(t, UserAgents, expectedUserAgents, "UserAgents slice does not match expected values")
}

func TestReferrers(t *testing.T) {
	// Verify the contents of the Referrers slice
	expectedReferrers := []string{
		"-", 
		"https://www.google.com", 
		"https://www.bing.com", 
		"https://www.example.com",
	}
	assert.Equal(t, Referrers, expectedReferrers, "Referrers slice does not match expected values")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


func TestSendResponse(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		success      bool
		message      string
		data         interface{}
		expectedBody string
		expectedCode int
	}{
		{
			name:         "Success with Data",
			statusCode:   http.StatusOK,
			success:      true,
			message:      "Request successful",
			data:         map[string]string{"key": "value"},
			expectedBody: `{"status":true,"message":"Request successful","data":{"key":"value"}}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Failure with Data",
			statusCode:   http.StatusBadRequest,
			success:      false,
			message:      "Bad Request",
			data:         map[string]string{"error": "invalid_input"},
			expectedBody: `{"status":false,"message":"Bad Request","data":{"error":"invalid_input"}}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Success without Data",
			statusCode:   http.StatusOK,
			success:      true,
			message:      "Request successful",
			data:         nil,
			expectedBody: `{"status":true,"message":"Request successful","data":null}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Failure without Data",
			statusCode:   http.StatusInternalServerError,
			success:      false,
			message:      "Internal Error",
			data:         nil,
			expectedBody: `{"status":false,"message":"Internal Error","data":null}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResponseRecorder (which is an in-memory ResponseWriter)
			rr := httptest.NewRecorder()

			// Initialize the ResponseHandler
			handler := &ResponseHandler{}

			// Call SendResponse
			handler.SendResponse(rr, tt.statusCode, tt.success, tt.message, tt.data)

			// Check if the status code matches
			assert.Equal(t, tt.expectedCode, rr.Code)

			// Check if the response body is as expected
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

// TestSendResponseError tests the error handling in SendResponse
func TestSendResponseError(t *testing.T) {

	// Set up the error scenario (invalid data causing marshaling failure)
	invalidData := make(chan int) // channels cannot be marshaled to JSON

	// Create a ResponseRecorder (which is an in-memory ResponseWriter)
	rr := httptest.NewRecorder()

	// Initialize the ResponseHandler
	handler := &ResponseHandler{}

	// Call SendResponse with invalid data
	handler.SendResponse(rr, http.StatusInternalServerError, false, "Internal Server Error", invalidData)

	// Check that the status code is 500
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check if the response body contains the appropriate error message

	exp_output := `Internal Server Error
`
	assert.Equal(t, exp_output, rr.Body.String())

}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestFirstLoad(t *testing.T){
	err := FirstLoad()
	assert.Equal(t, err, fmt.Errorf("error loading config from YAML: failed to read config.yaml: open config.yaml: no such file or directory"), "Error should not be there while loading from first load")
}

func TestGetEnvString(t *testing.T) {
	os.Setenv("key", "dummy")
	data1 := getEnvString("key", "default")

	assert.Equal(t, data1, "dummy", "Setting and unsetting of string environment variables")
	os.Unsetenv("key")

	data2 := getEnvString("dummy", "default")
	assert.Equal(t, data2, "default", "checking of string environment variable conversion")
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("key", strconv.Itoa(5))
	data1 := getEnvInt("key", 0)

	assert.Equal(t, data1, 5, "Setting and unsetting of integer environment variables")
	os.Unsetenv("key")

	data2 := getEnvInt("dummy", 0)
	assert.Equal(t, data2, 0, "Checking of integer environment variable conversion")
}

func TestLoadConfigFromYaml(t *testing.T) {
	fileData := []byte{}
	err := fmt.Errorf("open config.yaml: no such file or directory")
	expectedErr := fmt.Errorf("failed to read config.yaml: %v", err)
	actualErr := LoadConfigFromYaml(fileData, err)
	assert.Equal(t, expectedErr, actualErr, "Should return an error when file is not found")

	// Test 1: File Read Error (File not found)
	t.Run("File Read Error", func(t *testing.T) {
		err := fmt.Errorf("open config.yaml: no such file or directory")
		expectedErr := fmt.Errorf("failed to read config.yaml: %v", err)
		
		actualErr := LoadConfigFromYaml(nil, err)
		assert.Equal(t, expectedErr, actualErr, "Expected error when file can't be read")
	})

	// Test 2: Invalid YAML Format (YAML unmarshalling error)
	t.Run("Invalid YAML Format", func(t *testing.T) {
		// Simulate invalid YAML
		invalidYaml := []byte("{ invalid_yaml: ")

		// Simulate no read error (i.e., file is "read" but not valid)
		err := fmt.Errorf("yaml: line 1: did not find expected node content")
		expectedErr := fmt.Errorf("failed to parse config.yaml: %v", err)

		actualErr := LoadConfigFromYaml(invalidYaml, nil)
		assert.Equal(t, expectedErr, actualErr, "Expected error when YAML is invalid")
	})

	// Test 3: Successful YAML Parsing
	t.Run("Successful YAML Parsing", func(t *testing.T) {
		// Simulate valid YAML data
		validYaml := []byte(`
currentService:
  KEY_START_URL : "/logs"
  KEY_ALIVE_URL : "/"
  KEY_PORT : ":8080"

parserService:
#ENV PARSER_SERVICE_API="http://logparser:8082/logs"
  KEY_PARSER_API : "http://localhost:8083/logs"

#Current service configuration
KEY_RATE : 10
KEY_UNIT : "s"
`)

		// Simulate no error in reading
		err := LoadConfigFromYaml(validYaml, nil)

		// Assert the global data is correctly set
		assert.NoError(t, err, "Expected no error when YAML is valid")
		assert.Equal(t, ":8080", GloablMetaData.Port)
		assert.Equal(t, "/", GloablMetaData.IsAliveUrl)
		assert.Equal(t, "/logs", GloablMetaData.GenerateUrl)
		assert.Equal(t, "http://localhost:8083/logs", GloablMetaData.ProcessorApi)
		assert.Equal(t, int64(10), RateData.NumLogs)
		assert.Equal(t, "s", RateData.Unit)
	})

	// Test 4: Default Values (when RateData.NumLogs is 0 or Unit is invalid)
	t.Run("Default Values for Rate and Unit", func(t *testing.T) {
		// Case 1: No logs value in RateData
		RateData.NumLogs = 0
		RateData.Unit = "s"
		validYaml := []byte(`
currentService:
  KEY_START_URL : "/logs"
  KEY_ALIVE_URL : "/"
  KEY_PORT : ":8080"

parserService:
  KEY_PARSER_API : "http://localhost:8083/logs"

KEY_RATE : 15
KEY_UNIT : "s"
`)

		err := LoadConfigFromYaml(validYaml, nil)

		assert.NoError(t, err)
		assert.Equal(t, int64(15), RateData.NumLogs, "Expected NumLogs to be set from the config")
		assert.Equal(t, "s", RateData.Unit, "Expected Unit to be set from the config")

		// Case 2: Invalid unit in RateData
		RateData.Unit = "invalid"
		validYaml = []byte(`
currentService:
  KEY_START_URL : "/logs"
  KEY_ALIVE_URL : "/"
  KEY_PORT : ":8080"

parserService:
  KEY_PARSER_API : "http://localhost:8083/logs"

  KEY_RATE : 15
KEY_UNIT : "m"
`)

		err = LoadConfigFromYaml(validYaml, nil)

		assert.NoError(t, err)
		assert.Equal(t, "m", RateData.Unit, "Expected Unit to be set from the config")
	})
}

func TestReloadRateData(t *testing.T) {
	tests := []struct {
		name       string
		input      models.RequestPayload
		expectedErr bool
		expectedNumLogs int64
		expectedUnit string
	}{
		{
			name: "Valid input with positive rate and valid unit",
			input: models.RequestPayload{
				NumLogs: 10,
				Unit:    "s",
			},
			expectedErr: false,
			expectedNumLogs: 10,
			expectedUnit: "s",
		},
		{
			name: "Invalid NumLogs (less than or equal to zero)",
			input: models.RequestPayload{
				NumLogs: -5,
				Unit:    "s",
			},
			expectedErr: true,
			expectedNumLogs: 0,
			expectedUnit: "",
		},
		{
			name: "Invalid Unit (not s, m, or h)",
			input: models.RequestPayload{
				NumLogs: 10,
				Unit:    "invalid_unit",
			},
			expectedErr: true,
			expectedNumLogs: 0,
			expectedUnit: "",
		},
		{
			name: "Invalid NumLogs and Unit",
			input: models.RequestPayload{
				NumLogs: -10,
				Unit:    "invalid_unit",
			},
			expectedErr: true,
			expectedNumLogs: 0,
			expectedUnit: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global RateData before each test
			RateData = models.RequestPayload{}

			// Call ReloadRateData
			err := ReloadRateData(tt.input)

			// Assert if the error matches the expected value
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Assert the global RateData is updated correctly
			assert.Equal(t, tt.expectedNumLogs, RateData.NumLogs)
			assert.Equal(t, tt.expectedUnit, RateData.Unit)
		})
	}
}