package handlers

import (
	"LogParser/connection"
	"LogParser/handlers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
	handler := http.HandlerFunc(handlers.IsAlive)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("IsAlive returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := `{"status":true,"message":"Server 8084 is live\n","data":null}`
	actualResponse := rr.Body.String()
	assert.JSONEq(t, expectedResponse, actualResponse, "Response body doesn't match the expected format")

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
    handler := http.HandlerFunc(handlers.GetLogsCountHandler)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("GetLogsCountHandler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := `{"status":true,"message":"Logs count fetched successfully","data":{"count":5}}
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

    req, err := http.NewRequest("POST", "/addlogs", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.AddLogsHandler)
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
            "192.168.1.1", "-", "17/Mar/2025:13:30:20 +0530", "GET /home HTTP/1.1", 200,
            1234, "http://example.com", "Mozilla/5.0", "192.168.0.1",
        ),
    )
			
    req, err := http.NewRequest("GET", "/getlogs", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlers.GetLogsHandler)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("GetLogsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

	expected := `{"status":true,"message":"1 Logs fetched successfully","data":[{"remote_addr":"192.168.1.1","remote_user":"-","time_local":"17/Mar/2025:13:30:20 +0530","request":"GET /home HTTP/1.1","status":200,"body_bytes_sent":1234,"http_referer":"http://example.com","http_user_agent":"Mozilla/5.0","http_x_forwarded_for":"192.168.0.1"}]}`
    if rr.Body.String() != expected {
        t.Errorf("GetLogsHandler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unmet expectations: %s", err)
    }
}



