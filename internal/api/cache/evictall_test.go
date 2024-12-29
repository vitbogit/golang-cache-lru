package cache

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEvictAll_Success(t *testing.T) {
	// Create a new mock database
	mockService := new(MockService)

	// Set expectation
	mockService.On("EvictAll", context.Background()).Return(nil)

	// Create the handler with the mocked database
	handler := &Implementation{cacheService: mockService}

	// Create a new HTTP request to test the handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.EvictAll(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}

func TestEvictAll_UnknownError(t *testing.T) {
	// Create a new mock database
	mockService := new(MockService)

	// Set expectation
	mockService.On("EvictAll", mock.Anything).Return(fmt.Errorf("some error"))

	// Create the handler with the mocked database
	handler := &Implementation{cacheService: mockService}

	// Create a new HTTP request to test the handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.EvictAll(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}
