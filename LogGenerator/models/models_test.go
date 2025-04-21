package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestAllConfigModel(t *testing.T) {
	// Create a sample AllConfigModel with some values
	config := AllConfigModel{
		KEY_RATE: 100,
		KEY_UNIT: "minute",
		CurrentService: struct {
			KEY_START_URL string `yaml:"KEY_START_URL"`
			KEY_ALIVE_URL string `yaml:"KEY_ALIVE_URL"`
			KEY_PORT      string `yaml:"KEY_PORT"`
		}{
			KEY_START_URL: "/start",
			KEY_ALIVE_URL: "/alive",
			KEY_PORT:      "8080",
		},
		ParserService: struct {
			KEY_PARSER_API string `yaml:"KEY_PARSER_API"`
		}{
			KEY_PARSER_API: "http://localhost:5000/processLogs",
		},
	}

	// Marshaling the config struct to YAML
	marshalledYAML, err := yaml.Marshal(config)

	// We expect no error during marshalling
	assert.NoError(t, err, "Marshalling should not return an error")

	// Check the expected YAML output (the string should match the struct's YAML representation)
	expectedYAML := `
KEY_RATE: 100
KEY_UNIT: minute
currentService:
  KEY_START_URL: /start
  KEY_ALIVE_URL: /alive
  KEY_PORT: "8080"
parserService:
  KEY_PARSER_API: http://localhost:5000/processLogs
`
	assert.YAMLEq(t, expectedYAML, string(marshalledYAML), "The marshalled YAML should match the expected value")


	// Define the YAML configuration input string
	inputYAML := `
KEY_RATE: 100
KEY_UNIT: minute
currentService:
  KEY_START_URL: /start
  KEY_ALIVE_URL: /alive
  KEY_PORT: "8080"
parserService:
  KEY_PARSER_API: http://localhost:5000/processLogs
`
	// Create an empty AllConfigModel instance to unmarshal into
	config = AllConfigModel{}
	err = yaml.Unmarshal([]byte(inputYAML), &config)

	// We expect no error during unmarshalling
	assert.NoError(t, err, "Unmarshalling should not return an error")

	// Check if the unmarshalled struct has the expected values
	assert.Equal(t, 100, config.KEY_RATE, "The KEY_RATE should match the expected value")
	assert.Equal(t, "minute", config.KEY_UNIT, "The KEY_UNIT should match the expected value")
	assert.Equal(t, "/start", config.CurrentService.KEY_START_URL, "The KEY_START_URL should match the expected value")
	assert.Equal(t, "/alive", config.CurrentService.KEY_ALIVE_URL, "The KEY_ALIVE_URL should match the expected value")
	assert.Equal(t, "8080", config.CurrentService.KEY_PORT, "The KEY_PORT should match the expected value")
	assert.Equal(t, "http://localhost:5000/processLogs", config.ParserService.KEY_PARSER_API, "The KEY_PARSER_API should match the expected value")

	// Invalid YAML (missing KEY_RATE)
	invalidYAML := `
KEY_UNIT: minute
currentService:
  KEY_START_URL: /start
  KEY_ALIVE_URL: /alive
  KEY_PORT: "8080"
parserService:
  KEY_PARSER_API: http://localhost:5000/processLogs
`

	// Create an empty AllConfigModel instance to unmarshal into
	config = AllConfigModel{}
	err = yaml.Unmarshal([]byte(invalidYAML), &config)

	// We expect no error during unmarshalling, but the KEY_RATE should be zero (the default value)
	assert.NoError(t, err, "Unmarshalling should not return an error")

	// Check that the missing KEY_RATE field results in a zero value
	assert.Equal(t, 0, config.KEY_RATE, "The KEY_RATE should default to zero if not provided")

	// Empty YAML
	emptyYAML := ``
	// Create an empty AllConfigModel instance to unmarshal into
	config = AllConfigModel{}
	err = yaml.Unmarshal([]byte(emptyYAML), &config)

	// We expect no error during unmarshalling, but all fields should have default values
	assert.NoError(t, err, "Unmarshalling empty YAML should not return an error")

	// Check that all fields are set to their zero values (default values)
	assert.Equal(t, 0, config.KEY_RATE, "The KEY_RATE should default to zero if not provided")
	assert.Empty(t, config.KEY_UNIT, "The KEY_UNIT should be empty if not provided")
	assert.Empty(t, config.CurrentService.KEY_START_URL, "The KEY_START_URL should be empty if not provided")
	assert.Empty(t, config.CurrentService.KEY_ALIVE_URL, "The KEY_ALIVE_URL should be empty if not provided")
	assert.Empty(t, config.CurrentService.KEY_PORT, "The KEY_PORT should be empty if not provided")
	assert.Empty(t, config.ParserService.KEY_PARSER_API, "The KEY_PARSER_API should be empty if not provided")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestGlobalConstantvariablesMarshalling(t *testing.T) {
	// Create a sample GlobalConstantvariables with some values
	config := GlobalConstantvariables{
		Port:        "8080",
		IsAliveUrl:  "/alive",
		GenerateUrl: "/logs",
		ProcessorApi: "http://localhost:8082/logs",
	}

	// Marshaling the config struct to YAML
	marshalledYAML, err := yaml.Marshal(config)

	// We expect no error during marshalling
	assert.NoError(t, err, "Marshalling should not return an error")

	// Check the expected YAML output (the string should match the struct's YAML representation)
	expectedYAML := `
KEY_PORT: "8080"
KEY_ALIVE_URL: /alive
KEY_START_URL: /logs
KEY_PARSER_API: http://localhost:8082/logs
`
	assert.YAMLEq(t, expectedYAML, string(marshalledYAML), "The marshalled YAML should match the expected value")
}

func TestGlobalConstantvariablesUnmarshalling(t *testing.T) {
	// Define the YAML configuration input string
	inputYAML := `
KEY_PORT: "8080"
KEY_ALIVE_URL: /alive
KEY_START_URL: /logs
KEY_PARSER_API: http://localhost:8082/logs
`
	// Create an empty GlobalConstantvariables instance to unmarshal into
	var config GlobalConstantvariables
	err := yaml.Unmarshal([]byte(inputYAML), &config)

	// We expect no error during unmarshalling
	assert.NoError(t, err, "Unmarshalling should not return an error")

	// Check if the unmarshalled struct has the expected values
	assert.Equal(t, "8080", config.Port, "The Port should match the expected value")
	assert.Equal(t, "/alive", config.IsAliveUrl, "The IsAliveUrl should match the expected value")
	assert.Equal(t, "/logs", config.GenerateUrl, "The GenerateUrl should match the expected value")
	assert.Equal(t, "http://localhost:8082/logs", config.ProcessorApi, "The ProcessorApi should match the expected value")
}

func TestGlobalConstantvariablesUnmarshallingInvalidYAML(t *testing.T) {
	// Invalid YAML (missing KEY_PORT)
	invalidYAML := `
KEY_ALIVE_URL: /alive
KEY_START_URL: /logs
KEY_PARSER_API: http://localhost:8082/logs
`

	// Create an empty GlobalConstantvariables instance to unmarshal into
	var config GlobalConstantvariables
	err := yaml.Unmarshal([]byte(invalidYAML), &config)

	// We expect no error during unmarshalling, but the missing KEY_PORT should be empty or defaulted
	assert.NoError(t, err, "Unmarshalling should not return an error")

	// Check that the missing KEY_PORT field results in a zero value or an empty string
	assert.Equal(t, "", config.Port, "The Port should default to an empty string if not provided")
}

func TestGlobalConstantvariablesEmptyYAML(t *testing.T) {
	// Empty YAML
	emptyYAML := ``
	// Create an empty GlobalConstantvariables instance to unmarshal into
	var config GlobalConstantvariables
	err := yaml.Unmarshal([]byte(emptyYAML), &config)

	// We expect no error during unmarshalling, but all fields should have default values (empty strings)
	assert.NoError(t, err, "Unmarshalling empty YAML should not return an error")

	// Check that all fields are set to their zero values (default values)
	assert.Empty(t, config.Port, "The Port should be empty if not provided")
	assert.Empty(t, config.IsAliveUrl, "The IsAliveUrl should be empty if not provided")
	assert.Empty(t, config.GenerateUrl, "The GenerateUrl should be empty if not provided")
	assert.Empty(t, config.ProcessorApi, "The ProcessorApi should be empty if not provided")
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestRequestPayloadMarshalling(t *testing.T) {
	// Define a RequestPayload instance
	payload := RequestPayload{
		NumLogs: 1000,
		Unit:    "s",
	}

	// Marshal the RequestPayload struct to JSON
	jsonData, err := json.Marshal(payload)
	assert.NoError(t, err, "Marshalling should succeed")

	// Check the expected output of the marshalled JSON
	expectedJSON := `{"num_logs":1000,"time":"s"}`
	assert.JSONEq(t, expectedJSON, string(jsonData), "The marshalled JSON should match the expected output")
}

// TestRequestPayloadUnmarshalling tests the unmarshalling of JSON into the RequestPayload struct
func TestRequestPayloadUnmarshalling(t *testing.T) {
	jsonData := `{"num_logs":1000,"time":"s"}`
	var payload RequestPayload
	err := json.Unmarshal([]byte(jsonData), &payload)
	assert.NoError(t, err, "Unmarshalling should succeed")
	assert.Equal(t, int64(1000), payload.NumLogs)
	assert.Equal(t, "s", payload.Unit)
}

// TestRequestPayloadUnmarshallingInvalidJSON tests how the RequestPayload handles invalid or malformed JSON
func TestRequestPayloadUnmarshallingMissingNumLogsJSON(t *testing.T) {
	invalidJSON := `{"time":"s"}`

	var payload RequestPayload
	err := json.Unmarshal([]byte(invalidJSON), &payload)
	assert.NoError(t, err, "Default values should be set")
	
	assert.NoError(t, err, "Unmarshalling should succeed")
	assert.Equal(t, int64(0), payload.NumLogs)
	assert.Equal(t, "s", payload.Unit)

}

func TestRequestPayloadUnmarshallingMissingUnitJSON(t *testing.T) {
	invalidJSON := `{"num_logs": 1000}`

	var payload RequestPayload
	err := json.Unmarshal([]byte(invalidJSON), &payload)
	assert.NoError(t, err, "Default values should be set")
	
	assert.NoError(t, err, "Unmarshalling should succeed")
	assert.Equal(t, int64(1000), payload.NumLogs)
	assert.Equal(t, "", payload.Unit)

}

// TestRequestPayloadUnmarshallingEmptyJSON tests unmarshalling an empty JSON
func TestRequestPayloadUnmarshallingEmptyJSON(t *testing.T) {
	// Empty JSON
	emptyJSON := `{}`

	// Attempt to unmarshal the empty JSON into a RequestPayload struct
	var payload RequestPayload
	err := json.Unmarshal([]byte(emptyJSON), &payload)
	assert.NoError(t, err, "Default values should be set")
	// We expect an error because the fields are missing or have incorrect types
	//assert.Error(t, err, "Unmarshalling empty JSON should return an error")
	assert.Equal(t, int64(0), payload.NumLogs)
	assert.Equal(t, "", payload.Unit)
	//t.Log(payload)
}

// TestRequestPayloadEdgeCase tests the RequestPayload struct with edge case values
func TestRequestPayloadEdgeCase(t *testing.T) {
	// Edge case with zero logs and unit "0"
	payload := RequestPayload{
		NumLogs: 0,
		Unit:    "m",
	}

	// Marshal the RequestPayload struct to JSON
	jsonData, err := json.Marshal(payload)
	assert.NoError(t, err, "Marshalling should succeed")

	// Check the expected output of the marshalled JSON
	expectedJSON := `{"num_logs":0,"time":"m"}`
	assert.JSONEq(t, expectedJSON, string(jsonData), "The marshalled JSON should match the expected output")

	// Unmarshal the marshalled JSON back into a RequestPayload struct
	var unmarshalledPayload RequestPayload
	err = json.Unmarshal(jsonData, &unmarshalledPayload)
	assert.NoError(t, err, "Unmarshalling should succeed")

	// Validate that the edge case values are preserved correctly
	assert.Equal(t, int64(0), unmarshalledPayload.NumLogs)
	assert.Equal(t, "m", unmarshalledPayload.Unit)
}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////



func TestResponseMarshalling(t *testing.T) {
	successResponse := Response{
		Status:  true,
		Message: "Logs generated successfully",
		Data:    json.RawMessage(`[{"log": "data"}]`),
	}

	marshalledJSON, err := json.Marshal(successResponse)
	assert.NoError(t, err, "Marshalling should not return an error")

	expectedJSON := `{"status":true,"message":"Logs generated successfully","data":[{"log": "data"}]}`
	assert.JSONEq(t, expectedJSON, string(marshalledJSON), "The marshalled JSON should match the expected value")
	failedResponse := Response{
		Status:  false,
		Message: "Failed to generate logs",
		Data:    nil,
	}
	marshalledFailedJSON, err := json.Marshal(failedResponse)
	assert.NoError(t, err, "Marshalling should not return an error")

	expectedFailedJSON := `{"status":false,"message":"Failed to generate logs","data":null}`
	assert.JSONEq(t, expectedFailedJSON, string(marshalledFailedJSON), "The marshalled failed response JSON should match the expected value")
}

func TestResponseUnmarshalling(t *testing.T) {
	validJSON := `{"status":true,"message":"Logs generated successfully","data":[{"log": "data"}]}`

	var validResponse Response
	err := json.Unmarshal([]byte(validJSON), &validResponse)
	assert.NoError(t, err, "Unmarshalling valid JSON should not return an error")
	assert.True(t, validResponse.Status, "The status should be true")
	assert.Equal(t, "Logs generated successfully", validResponse.Message, "The message should match the expected string")
	assert.JSONEq(t, `[{"log": "data"}]`, string(validResponse.Data), "The data should match the expected JSON")
	
	invalidJSON := `{"status":true,"message":"Missing data"}`
	var invalidResponse Response
	err = json.Unmarshal([]byte(invalidJSON), &invalidResponse)
	assert.NoError(t, err, "Unmarshalling invalid JSON should not return an error")
	assert.Nil(t, invalidResponse.Data, "Data should be nil if it's missing in the JSON")
}

func TestResponseValidation(t *testing.T) {
	// Test the marshalled result when fields are empty
	emptyResponse := Response{
		Status:  false,
		Message: "",
		Data:    nil,
	}

	// Marshalling the empty response to JSON
	marshalledEmptyJSON, err := json.Marshal(emptyResponse)

	// We expect no error during marshalling
	assert.NoError(t, err, "Marshalling empty response should not return an error")

	// Check if the JSON contains the expected empty fields
	expectedEmptyJSON := `{"status":false,"message":"","data":null}`
	assert.JSONEq(t, expectedEmptyJSON, string(marshalledEmptyJSON), "The marshalled empty response JSON should match the expected value")
}
