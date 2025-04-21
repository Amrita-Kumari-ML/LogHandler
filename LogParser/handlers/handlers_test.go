package handlers

import (
	"LogParser/connection"
	"LogParser/logger"
	"LogParser/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestIsAlive(t *testing.T) {
	//connection.InitDB()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IsAlive)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("IsAlive returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := `{"status":true,"message":"Server  is live","data":null}`
	actualResponse := rr.Body.String()
	assert.JSONEq(t, expectedResponse, actualResponse, "Response body doesn't match the expected format")

}



func init() {
	logger.InitLogger("error") // suppress debug/info in tests
}

// Mock versions of the handlers for testing call routing
var getCalled, postCalled, deleteCalled bool

func TestHandleType(t *testing.T) {

	tests := []struct {
		method        string
		expectedCode  int
		expectedMsg   string
		expectGet     bool
		expectPost    bool
		expectDelete  bool
	}{
		{"GET", http.StatusOK, "Mock Get Called", true, false, false},
		{"POST", http.StatusOK, "Mock Post Called", false, true, false},
		{"DELETE", http.StatusOK, "Mock Delete Called", false, false, true},
		{"PUT", http.StatusMethodNotAllowed, "Only GET, POST, DELETE methods are allowed to execute the task", false, false, false},
	}

			req := httptest.NewRequest(tests[3].method, "/logs", nil)
			rr := httptest.NewRecorder()

			HandleType(rr, req)

			resp := rr.Result()
			assert.Equal(t, tests[3].expectedCode, resp.StatusCode)

			body := rr.Body.String()
			assert.Contains(t, body, tests[3].expectedMsg)

			assert.Equal(t, tests[3].expectGet, getCalled)
			assert.Equal(t, tests[3].expectPost, postCalled)
			assert.Equal(t, tests[3].expectDelete, deleteCalled)

	
}

func TestGetLogsCountHandler_DBConnectionFail(t *testing.T) {
	// Simulate DB connection failure
	connection.DB = nil

	req := httptest.NewRequest("GET", "/logs/count", nil)
	rr := httptest.NewRecorder()

	GetLogsCountHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), `"status":false`)
	assert.Contains(t, rr.Body.String(), `Failed to connect to Database`)
}

func TestFormatTime_WithValidTime(t *testing.T) {
	// Create a known time
	inputTime := time.Date(2025, time.April, 10, 15, 30, 0, 0, time.UTC)
	expected := inputTime.Format(time.RFC3339)

	// Call formatTime
	result := FormatTime(&inputTime)

	assert.NotNil(t, result)
	assert.Equal(t, expected, *result)
}

func TestFormatTime_WithNil(t *testing.T) {
	var tNil *time.Time

	result := FormatTime(tNil)

	assert.Nil(t, result)
}




func TestGetLogsCountHandler(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock database: %s", err)
    }
    defer db.Close()
    mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM logs WHERE 1=1").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
    connection.DB = db
    req, err := http.NewRequest("GET", "/getlogsCount?remote_addr=127.0.0.1", nil) 
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetLogsCountHandler)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("GetLogsCountHandler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := `{"status":true,"message":"Logs Found Success","data":{"fetch":5,"total":0}}
`
    if rr.Body.String() != expected {
        t.Errorf("GetLogsCountHandler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }

}


// Test for AddLogsHandler with mock database
func TestAddLogsHandler(t *testing.T) {
    // Mocking database
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock database: %s", err)
    }
    defer db.Close()

	connection.DB = db
    mock.ExpectExec("INSERT INTO logs").WillReturnResult(sqlmock.NewResult(1, 1))
    logs := []string{
        "192.168.1.1 - - [17/Mar/2025:13:30:20 +0530] \"GET /home HTTP/1.1\" 200 1180 \"https://www.bing.com\" \"Mozilla/5.0...\"",
    }
    jsonStr, err := json.Marshal(logs)
    if err != nil {
        t.Fatalf("Failed to marshal logs: %v", err)
    }

    req, err := http.NewRequest("POST", "/logs", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(AddLogsHandler)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("AddLogsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := `{"status":true,"message":"Logs stored successfully, 1 rows inserted.","data":null}
`
    if rr.Body.String() != expected {
        t.Errorf("AddLogsHandler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}


func TestGetLogsHandler(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock database: %s", err)
    }
    defer db.Close()

    connection.DB = db
	mock.ExpectQuery("SELECT remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for").
    WillReturnRows(
        sqlmock.NewRows([]string{
            "remote_addr", "remote_user", "time_local", "request", "status",
            "body_bytes_sent", "http_referer", "http_user_agent", "http_x_forwarded_for",
        }).AddRow(
            "192.168.1.1", "-",
            time.Date(2025, time.March, 17, 13, 30, 20, 0, time.FixedZone("IST", 19800)), // ✅ FIXED here
            "GET /home HTTP/1.1", 200,
            1234, "http://example.com", "Mozilla/5.0", "192.168.0.1",
        ),
    )
			
    req, err := http.NewRequest("GET", "/logs", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetLogsHandler)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("GetLogsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

	expected := `{"status":true,"message":"Fetched logs successfully","data":{"count":{"fetch":1,"total":0},"logs":[{"remote_addr":"192.168.1.1","remote_user":"-","time_local":"2025-03-17T13:30:20+05:30","request":"GET /home HTTP/1.1","status":200,"body_bytes_sent":1234,"http_referer":"http://example.com","http_user_agent":"Mozilla/5.0","http_x_forwarded_for":"192.168.0.1"}],"paging":{"limit":10,"next_cursor":null,"prev_cursor":"2025-03-17T13:30:20+05:30"}}}
`
    if rr.Body.String() != expected {
        t.Errorf("GetLogsHandler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unmet expectations: %s", err)
    }
}
	


func TestInsertOneLog_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	connection.DB = db // Set mock DB

	log := models.Log{
		RemoteAddr:        "127.0.0.1",
		RemoteUser:        "-",
		TimeLocal:         time.Now(),
		Request:           "GET /home HTTP/1.1",
		Status:            200,
		BodyBytesSent:     500,
		HttpReferer:       "http://example.com",
		HttpUserAgent:     "Mozilla/5.0",
		HttpXForwardedFor: "192.168.0.1",
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(log.RemoteAddr, log.RemoteUser, log.TimeLocal, log.Request, log.Status, log.BodyBytesSent, log.HttpReferer, log.HttpUserAgent, log.HttpXForwardedFor).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = InsertOneLog(log)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertOneLog_DBDown(t *testing.T) {
	connection.DB = nil // Simulate DB not alive

	log := models.Log{}
	err := InsertOneLog(log)
	assert.Error(t, err)
	assert.Equal(t, "Database is down!", err.Error())
}

func TestInsertOneLog_InsertFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	connection.DB = db

	log := models.Log{}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(log.RemoteAddr, log.RemoteUser, log.TimeLocal, log.Request, log.Status, log.BodyBytesSent, log.HttpReferer, log.HttpUserAgent, log.HttpXForwardedFor).
		WillReturnError(assert.AnError)

	err = InsertOneLog(log)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProcessLogWorker(t *testing.T) {
	logs := make(chan string, 1)
	results := make(chan models.Log, 1)
	var wg sync.WaitGroup

	// Add one item to WaitGroup as one goroutine will run
	wg.Add(1)
	go ProcessLogWorker(logs, results, &wg)

	// Send a test log line
	logs <- `127.0.0.1 - - [17/Mar/2025:13:30:20 +0530] "GET /home HTTP/1.1" 200 500 "-" "Mozilla/5.0" "192.168.0.1"`
	close(logs) // Important to close channel so goroutine can exit

	// Wait for goroutine to finish
	wg.Wait()
	close(results)

	// Assert the result
	parsedLog := <-results
	assert.Equal(t, "127.0.0.1", parsedLog.RemoteAddr)
	assert.Equal(t, "GET /home HTTP/1.1", parsedLog.Request)
	assert.Equal(t, 200, parsedLog.Status)
}

func TestParseLog_Valid(t *testing.T) {
	logLine := `192.168.1.1 - user123 [2025-04-10T10:20:30Z] "GET /api HTTP/1.1" 200 512 "http://example.com" "Go-http-client/1.1" "192.168.1.100"`

	log := ParseLog(logLine)

	assert.Equal(t, "192.168.1.1", log.RemoteAddr)
	assert.Equal(t, "user123", log.RemoteUser)
	assert.Equal(t, "GET /api HTTP/1.1", log.Request)
	assert.Equal(t, 200, log.Status)
	assert.Equal(t, 512, log.BodyBytesSent)
	assert.Equal(t, "http://example.com", log.HttpReferer)
	assert.Equal(t, "Go-http-client/1.1", log.HttpUserAgent)
	assert.Equal(t, "192.168.1.100", log.HttpXForwardedFor)
	assert.Equal(t, time.Date(2025, 4, 10, 10, 20, 30, 0, time.UTC), log.TimeLocal)
}

func TestParseLog_InvalidFormat(t *testing.T) {
	logLine := `This is a malformed log line`
	log := ParseLog(logLine)

	assert.Equal(t, models.Log{}, log)
}

func TestParseLog_InvalidTime(t *testing.T) {
	logLine := `192.168.1.1 - user123 [invalid-time-format] "GET /api HTTP/1.1" 200 512 "http://example.com" "Go-http-client/1.1" "192.168.1.100"`
	log := ParseLog(logLine)

	assert.Equal(t, time.Time{}, log.TimeLocal) // should be zero time
	assert.Equal(t, "192.168.1.1", log.RemoteAddr)
}

func TestAtoi_ValidInput(t *testing.T) {
	assert.Equal(t, 123, Atoi("123"))
	assert.Equal(t, 0, Atoi("0"))
	assert.Equal(t, -42, Atoi("-42"))
}

func TestAtoi_InvalidInput(t *testing.T) {
	// Should return 0 for invalid input as per current implementation
	assert.Equal(t, 0, Atoi("abc"))
	assert.Equal(t, 0, Atoi(""))
	assert.Equal(t, 0, Atoi("12a3"))
}

/*
// TestGetLogsHandler tests the GetLogsHandler function
func TestGetLogsHandler(t *testing.T) {
	// Set up mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// Mock database query and expected return values
	mock.ExpectQuery(utils.QUERY_COUNT_ALL).
		WillReturnRows(sqlmock.NewRows([]string{"total_logs"}).AddRow(10))

	mock.ExpectQuery("SELECT remote_addr, remote_user, time_local, request, status, body_bytes_sent, http_referer, http_user_agent, http_x_forwarded_for").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"remote_addr", "remote_user", "time_local", "request", "status",
				"body_bytes_sent", "http_referer", "http_user_agent", "http_x_forwarded_for",
			}).AddRow(
				"192.168.1.1", "-", "17/Mar/2025:13:30:20 +0530", "GET /home HTTP/1.1", 200,
				1234, "http://example.com", "Mozilla/5.0", "192.168.0.1",
			),
		)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the GetLogsHandler
	handler := http.HandlerFunc(GetLogsHandler)
	handler.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, 500, rr.Code)

	// Expected JSON response structure
	expectedResponse := `{"status":true,"message":"Fetched logs successfully","data":{"count":{"total":10,"fetch":1},"logs":[{"remote_addr":"192.168.1.1","remote_user":"-","time_local":"17/Mar/2025:13:30:20 +0530","request":"GET /home HTTP/1.1","status":200,"body_bytes_sent":1234,"http_referer":"http://example.com","http_user_agent":"Mozilla/5.0","http_x_forwarded_for":"192.168.0.1"}],"paging":{"next_cursor":null,"prev_cursor":null,"limit":10}}}`

	// Assert response body
	assert.JSONEq(t, expectedResponse, rr.Body.String())

	// Ensure all expectations were met with the mock database
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %s", err)
	}
}

// TestGetLogsHandler_DBError tests the scenario when the database is not available
func TestGetLogsHandler_DBError(t *testing.T) {
	// Set up mock database connection
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the GetLogsHandler
	handler := http.HandlerFunc(GetLogsHandler)
	handler.ServeHTTP(rr, req)

	// Assert status code and error message when DB is down
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to connect to Database!")
}

// TestGetLogsHandler_QueryError tests the scenario when there's an error in fetching logs from the database
func TestGetLogsHandler_QueryError(t *testing.T) {
	// Set up mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// Mock the query to return an error
	mock.ExpectQuery(utils.QUERY_COUNT_ALL).WillReturnError(fmt.Errorf("failed to fetch total logs"))

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the GetLogsHandler
	handler := http.HandlerFunc(GetLogsHandler)
	handler.ServeHTTP(rr, req)

	// Assert status code and error message when the query fails
	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to query database")
}
	*/