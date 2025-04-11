package models

import (
	_"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSendResponse_WithData(t *testing.T) {
	rr := httptest.NewRecorder()

	mockData := map[string]string{"key": "value"}
	SendResponse(rr, http.StatusOK, true, "Success", mockData)

	result := rr.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

	bodyBytes, _ := io.ReadAll(result.Body)

	var responseBody map[string]interface{}
	err := json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, true, responseBody["status"])
	assert.Equal(t, "Success", responseBody["message"])
	assert.NotNil(t, responseBody["data"])
}

func TestSendResponse_WithoutData(t *testing.T) {
	rr := httptest.NewRecorder()

	SendResponse(rr, http.StatusOK, true, "No data", nil)

	result := rr.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

	bodyBytes, _ := io.ReadAll(result.Body)

	var responseBody map[string]interface{}
	err := json.Unmarshal(bodyBytes, &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, true, responseBody["status"])
	assert.Equal(t, "No data", responseBody["message"])
	assert.Nil(t, responseBody["data"])
}

func TestSendResponse_MarshalError(t *testing.T) {
	rr := httptest.NewRecorder()

	// Channels cannot be marshalled to JSON; this will trigger a marshal error
	badData := make(chan int)

	SendResponse(rr, http.StatusOK, true, "Invalid data", badData)

	result := rr.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	bodyBytes, _ := io.ReadAll(result.Body)
	assert.Equal(t, "Internal Server Error\n", string(bodyBytes))
}
