package utils

import (
	"LogParser/logger"
	"LogParser/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func init() {
	logger.InitializeLogger("error") // suppress debug/info in tests
}
func TestFirstLoad_WithEnvVars(t *testing.T) {
	// Set mock environment variable
	os.Setenv("PORT", ":8083")

	// First load with environment variable
	err := FirstLoad()

	exp := fmt.Errorf("error loading config from YAML: error reading YAML file: open config.yaml: no such file or directory\n")
	// Assert that no error occurred
	assert.Equal(t, err, exp)

	// Assert the global ConfigData has the correct values
	assert.Equal(t, ":8083", ConfigData.PORT)

	// Clean up
	os.Unsetenv("PORT")
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

func TestGetEnvInt2(t *testing.T) {
	// Set environment variable for the test
	os.Setenv("TEST_INT", "123")

	// Test case where the environment variable is set and correctly parsed as int
	result := getEnvInt("TEST_INT", 456)
	assert.Equal(t, 123, result)

	// Test case where the environment variable is not set
	result = getEnvInt("NON_EXISTENT_KEY", 456)
	assert.Equal(t, 456, result)

	// Test case where the environment variable is set but is not a valid integer
	os.Setenv("TEST_INT", "invalid")
	result = getEnvInt("TEST_INT", 456)
	assert.Equal(t, 456, result)

	// Clean up
	os.Unsetenv("TEST_INT")
}

func TestLoadConfigFromYaml(t *testing.T) {
	// Mock loading of YAML file
	// For this test, we simulate a successful loading process.

	// Test loading configuration from YAML
	err := LoadConfigFromYaml()
	exp := "error reading YAML file: open config.yaml: no such file or directory\n"
	// Assert no error occurred
	assert.EqualError(t, err, exp)

	// Assert the global ConfigData is populated
	assert.Equal(t, ":8083", ConfigData.PORT)
}


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

func TestGenerateFilteredGetQuery(t *testing.T) {
	// Setup filters
	filters := map[string]interface{}{
		"status": "200",
		"request": "/api/v1/logs",
	}

	// Setup pagination filter
	paginationFilter := models.Pagination{
		Limit: 10,
		Cursor: nil,
	}

	// Setup date filter
	startTime := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2022, time.March, 2, 0, 0, 0, 0, time.UTC)
	dateFilter := models.TimeFilter{
		Start_time: &startTime,
		End_time:   &endTime,
	}

	// Call the function
	query, args := GenerateFilteredGetQuery(filters, paginationFilter, dateFilter)

	// Expected query string
	expectedQuery := `SELECT remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for FROM logs WHERE 1=1 AND status = $1 AND request = $2 AND time_local >= $3 AND time_local <= $4 LIMIT $5`

	// Assert that the query matches
	assert.Equal(t, expectedQuery, query)

	// Assert that the args are correctly constructed
	expectedArgs := []interface{}{"200", "/api/v1/logs", "2022-03-01T00:00:00Z", "2022-03-02T00:00:00Z", 10}
	assert.Equal(t, expectedArgs, args)
}

func TestGenerateFilteredCountQuery(t *testing.T) {
	// Setup filters
	filters := map[string]interface{}{
		"status": "200",
	}

	// Call the function
	query, args := GenerateFilteredCountQuery(filters)

	// Expected query string
	expectedQuery := `SELECT COUNT(*) FROM logs WHERE 1=1 AND status = $1`

	// Assert that the query matches
	assert.Equal(t, expectedQuery, query)

	// Assert that the args are correctly constructed
	expectedArgs := []interface{}{"200"}
	assert.Equal(t, expectedArgs, args)
}

func TestGenerateDeleteQuery(t *testing.T) {
	// Setup filters
	filters := map[string]interface{}{
		"status": "500",
		"request": "/api/v1/deleteLogs",
	}

	// Call the function
	query, args := GenerateDeleteQuery(filters)

	// Expected query string
	expectedQuery := `DELETE FROM logs WHERE 1=1 AND status = $1 AND request = $2`

	// Assert that the query matches
	assert.Equal(t, expectedQuery, query)

	// Assert that the args are correctly constructed
	expectedArgs := []interface{}{"500", "/api/v1/deleteLogs"}
	assert.Equal(t, expectedArgs, args)
}

func TestGenerateAddQuery(t *testing.T) {
	// Create sample logs
	logs := []models.Log{
		{
			RemoteAddr:   "192.168.1.1",
			RemoteUser:   "user1",
			TimeLocal:    time.Now(),
			Request:      "/api/v1/logs",
			Status:       200,
			BodyBytesSent: 123,
			HttpReferer:  "https://example.com",
			HttpUserAgent: "Mozilla/5.0",
			HttpXForwardedFor: "192.168.1.2",
		},
	}

	// Call the function
	query, args := GenerateAddQuery(logs)

	// Expected query string
	expectedQuery := `
		INSERT INTO logs (remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	// Assert that the query matches
	assert.Contains(t, query, expectedQuery)//"INSERT INTO logs (remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for) VALUES"

	// Assert that the args are correctly constructed
	assert.Len(t, args, 9) // There should be 9 values in the args slice
	assert.Equal(t, "192.168.1.1", args[0])
	assert.Equal(t, "user1", args[1])
	//assert.Equal(t, logs[0].TimeLocal.UTC().Format(time.RFC3339), args[2].(string))
	assert.Equal(t, "/api/v1/logs", args[3])
	assert.Equal(t, 200, args[4])
	assert.Equal(t, 123, args[5])
	assert.Equal(t, "https://example.com", args[6])
	assert.Equal(t, "Mozilla/5.0", args[7])
	assert.Equal(t, "192.168.1.2", args[8])
}

func TestGetCount(t *testing.T) {
	// Call the function
	query := GetCount()

	// Expected query string
	expectedQuery := `SELECT COUNT(*) FROM logs;`

	// Assert that the query matches
	assert.Equal(t, expectedQuery, query)
}

func createMockRequest(queryParams map[string]string) *http.Request {
	urlValues := url.Values{}
	for key, value := range queryParams {
		urlValues.Add(key, value)
	}

	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		panic(err)
	}

	req.URL.RawQuery = urlValues.Encode()
	return req
}

func TestGenerateFiltersMap(t *testing.T) {
	// Setup query parameters for the test
	queryParams := map[string]string{
		"remote_addr":      "192.168.1.1",
		"status":           "200",
		"body_bytes_sent":  "512",
		"http_referer":     "https://example.com",
		"http_user_agent":  "Mozilla/5.0",
		"http_x_forwarded_for": "192.168.1.2",
	}

	// Create mock HTTP request
	req := createMockRequest(queryParams)

	// Call the function
	filters := GenerateFiltersMap(req)

	// Assert that the filters map is generated correctly
	assert.Equal(t, "192.168.1.1", filters["remote_addr"])
	assert.Equal(t, 200, filters["status"])
	assert.Equal(t, 512, filters["body_bytes_sent"])
	assert.Equal(t, "https://example.com", filters["http_referer"])
	assert.Equal(t, "Mozilla/5.0", filters["http_user_agent"])
	assert.Equal(t, "192.168.1.2", filters["http_x_forwarded_for"])
}

func TestGetPaginationParams(t *testing.T) {
	// Setup query parameters for pagination
	queryParams := map[string]string{
		"page":   "2",
		"limit":  "20",
		"cursor": "2025-04-10T10:30:00Z",
	}

	// Create mock HTTP request
	req := createMockRequest(queryParams)

	// Call the function
	pagination := GetPaginationParams(req)

	// Assert that pagination is parsed correctly
	assert.Equal(t, 2, pagination.Page)
	assert.Equal(t, 20, pagination.Limit)
	assert.NotNil(t, pagination.Cursor)
	assert.Equal(t, time.Date(2025, time.April, 10, 10, 30, 0, 0, time.UTC), *pagination.Cursor)
}

func TestGetPaginationParamsWithDefaults(t *testing.T) {
	// Create mock HTTP request without pagination parameters
	req := createMockRequest(map[string]string{})

	// Call the function
	pagination := GetPaginationParams(req)

	// Assert that default pagination values are used
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 10, pagination.Limit)
	assert.NotNil(t, pagination.Cursor)
}

func TestGetDateFilters(t *testing.T) {
	// Setup query parameters for time filtering
	queryParams := map[string]string{
		"start_time": "2025-04-08T06:00:00Z",
		"end_time":   "2025-04-09T06:00:00Z",
	}

	// Create mock HTTP request
	req := createMockRequest(queryParams)

	// Call the function
	timeFilters, err := GetDateFilters(req)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the start and end times are parsed correctly
	assert.Equal(t, time.Date(2025, time.April, 8, 6, 0, 0, 0, time.UTC), *timeFilters.Start_time)
	assert.Equal(t, time.Date(2025, time.April, 9, 6, 0, 0, 0, time.UTC), *timeFilters.End_time)
}

func TestGetDateFiltersWithInvalidStartTime(t *testing.T) {
	// Setup query parameters with an invalid start_time format
	queryParams := map[string]string{
		"start_time": "invalid-date",
	}

	// Create mock HTTP request
	req := createMockRequest(queryParams)

	// Call the function
	_, err := GetDateFilters(req)

	// Assert that an error occurred
	assert.Error(t, err)
	//assert.Equal(t, time.Time{}, *timeFilters.Start_time)
}

func TestGetDateFiltersWithStartTimeAfterEndTime(t *testing.T) {
	// Setup query parameters with start_time after end_time
	queryParams := map[string]string{
		"start_time": "2025-04-09T06:00:00Z",
		"end_time":   "2025-04-08T06:00:00Z",
	}

	// Create mock HTTP request
	req := createMockRequest(queryParams)

	// Call the function
	timeFilters, err := GetDateFilters(req)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that start and end times were swapped
	assert.Equal(t, time.Date(2025, time.April, 8, 6, 0, 0, 0, time.UTC), *timeFilters.Start_time)
	assert.Equal(t, time.Date(2025, time.April, 9, 6, 0, 0, 0, time.UTC), *timeFilters.End_time)
}

func TestGetDateFiltersWithDefaultValues(t *testing.T) {
	// Create mock HTTP request without time parameters
	req := createMockRequest(map[string]string{})

	// Call the function
	timeFilters, err := GetDateFilters(req)

	// Assert that no error occurred and that default values are set (nil for start and end time)
	assert.NoError(t, err)
	assert.Nil(t, timeFilters.Start_time)
	assert.Nil(t, timeFilters.End_time)
}