package server

import (
	"LogGenerator/loggenerator"
	"LogGenerator/logger"
	"LogGenerator/models"
	"LogGenerator/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

 var yaml = []byte(`
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

func TestIsAlive(t *testing.T) {
	logger.InitializeLogger("info")
	utils.LoadConfigFromYaml(yaml, nil)
	handler := &ServerHandler{
		ResponseW: &utils.ResponseHandler{},
		LogGen:    nil,
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.IsAlive(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("IsAlive returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "{\"status\":true,\"message\":\"Server :8080 is live\",\"data\":null}\n"
	if rr.Body.String() != expected {
		t.Errorf("IsAlive returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}


func TestLogTestHandler_ValidRequest(t *testing.T) {
	logger.InitializeLogger("debug")
	utils.LoadConfigFromYaml(yaml, nil)
	handler := &ServerHandler{
		ResponseW: &utils.ResponseHandler{},
		LogGen: &loggenerator.Generator{},
	}
	rateModel := models.RequestPayload{
		NumLogs: 2,
		Unit: "s",
	}

	payload, err := json.Marshal(rateModel)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/gen", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler.LogHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
	}

	expected := "{\"status\":true,\"message\":\"Task is in progress...\",\"data\":null}\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %v, but got %v", expected, rr.Body.String())
	}
}


func TestLogTestHandler_InvalidMethod(t *testing.T) {
	logger.InitializeLogger("debug")
	utils.LoadConfigFromYaml(yaml, nil)
	req, err := http.NewRequest(http.MethodGet, "/gen", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	serv := &ServerHandler{
		ResponseW: &utils.ResponseHandler{},
		LogGen: &loggenerator.Generator{},
	}

	serv.LogHandler(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %v, but got %v", http.StatusMethodNotAllowed, status)
	}
	expected := "{\"status\":false,\"message\":\"Only POST method allowed\",\"data\":null}\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %v, but got %v", expected, rr.Body.String())
	}
}


func TestLogTestHandler_MissingUnit(t *testing.T) {
	logger.InitializeLogger("debug")
	utils.LoadConfigFromYaml(yaml, nil)
	rateModel := models.RequestPayload{
		NumLogs: 10,
	}

	payload, err := json.Marshal(rateModel)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/gen", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	serv := &ServerHandler{
		ResponseW: &utils.ResponseHandler{},
		LogGen:    &loggenerator.Generator{},
	}

	serv.LogHandler(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %v, but got %v", http.StatusBadRequest, status)
	}

	expected := "{\"status\":false,\"message\":\"Invalid unit. Use s, m, or h for unit variable\",\"data\":null}\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %v, but got %v", expected, rr.Body.String())
	}
}


func TestLogTestHandler_InvalidUnit(t *testing.T) {
	logger.InitializeLogger("debug")
	utils.LoadConfigFromYaml(yaml, nil)
	rateModel := models.RequestPayload{
		NumLogs: 10,
		Unit: "xyz",
	}

	payload, err := json.Marshal(rateModel)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/gen", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	serv := &ServerHandler{
		ResponseW: &utils.ResponseHandler{},
		LogGen:    &loggenerator.Generator{},
	}

	serv.LogHandler(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %v, but got %v", http.StatusBadRequest, status)
	}

	expected := "{\"status\":false,\"message\":\"Invalid unit. Use s, m, or h for unit variable\",\"data\":null}\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %v, but got %v", expected, rr.Body.String())
	}
}