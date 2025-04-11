package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerHandlerInitialization(t *testing.T) {
	sHandler := &ServerHandler{}
	// Assert that the ServerHandler is initialized correctly
	assert.NotNil(t, sHandler)
}

func TestHandlerMapperInitialization(t *testing.T) {
	serverHandler := &ServerHandler{}
	// Create HandlerMapper instance using NewHandlerMapper
	handlerMapper := NewHandlerMapper(serverHandler)

	// Assert that the HandlerMapper is initialized correctly
	assert.NotNil(t, handlerMapper)
	assert.Equal(t, serverHandler, handlerMapper.ServerHandler)
}


