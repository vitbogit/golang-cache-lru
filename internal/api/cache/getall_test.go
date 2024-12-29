package cache

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll_NoContent(t *testing.T) {
	// Create a new mock service
	mockService := new(MockService)

	// Set expectation
	mockService.On("GetAll", context.Background()).Return(nil, nil, nil)

	// Create the handler with the mocked service
	handler := &Implementation{cacheService: mockService}

	// Create a new HTTP request to test the handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.GetAll(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}

func TestGetAll_OK(t *testing.T) {
	// Create a new mock service
	mockService := new(MockService)

	// Set expectation
	mockService.On("GetAll", context.Background()).Return([]string{"key1", "key2"}, []interface{}{"key1", 2}, nil)

	// Create the handler with the mocked service
	handler := &Implementation{cacheService: mockService}

	// Create a new HTTP request to test the handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.GetAll(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	// TODO: check for returning values in JSON

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}

func TestGetAll_Error(t *testing.T) {
	// Create a new mock service
	mockService := new(MockService)

	// Set expectation
	mockService.On("GetAll", context.Background()).Return(nil, nil, fmt.Errorf("some error"))

	// Create the handler with the mocked service
	handler := &Implementation{cacheService: mockService}

	// Create a new HTTP request to test the handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.GetAll(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}
