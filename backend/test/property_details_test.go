package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"interiorshuffle.com/property_app/api"
)

// MockRedisClient is a mock implementation of the RedisClient interface
type MockRedisClient struct {
	mock.Mock
}

// Ensure MockRedisClient implements the RedisClient interface
var _ api.RedisClient = (*MockRedisClient)(nil)

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

// MockHTTPClient mocks the HTTP client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetPropertyDetail(t *testing.T) {
	// Setup
	mockRedis := new(MockRedisClient)
	mockHTTP := new(MockHTTPClient)

	// Create a property detail object
	expectedProperty := &api.PropertyDetail{
		ID:      "property123",
		Address: "123 Main St",
		City:    "Arlington",
		State:   "VA",
		Zip:     "22205",
		Price:   "$500,000",
	}

	// Convert to JSON for consistent formatting
	mockAPIResponseBytes, _ := json.Marshal(expectedProperty)
	mockAPIResponse := string(mockAPIResponseBytes)

	// Mock Redis Get behavior (simulate cache miss)
	mockRedis.On("Get", mock.Anything, "property123").Return(redis.NewStringResult("", redis.Nil))

	// Mock Redis Set behavior using the same JSON structure
	mockRedis.On("Set", mock.Anything, "property123", mockAPIResponse, time.Hour).Return(redis.NewStatusResult("OK", nil))

	// Mock HTTP response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(mockAPIResponse)),
	}
	mockHTTP.On("Do", mock.Anything).Return(mockResp, nil)

	// Replace default HTTP client with mock
	originalHTTPClient := api.DefaultHTTPClient
	api.DefaultHTTPClient = mockHTTP
	defer func() { api.DefaultHTTPClient = originalHTTPClient }()

	// Call the function
	result, err := api.GetPropertyDetail(mockRedis, nil, "property123", "user123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedProperty, result)

	// Verify all expectations were met
	mockRedis.AssertExpectations(t)
	mockHTTP.AssertExpectations(t)
}

func TestGetPropertyDetailCacheHit(t *testing.T) {
	mockRedis := new(MockRedisClient)

	// Create a property detail object
	expectedProperty := &api.PropertyDetail{
		ID:      "property123",
		Address: "123 Main St",
		City:    "Arlington",
		State:   "VA",
		Zip:     "22205",
		Price:   "$500,000",
	}

	// Convert to JSON for consistent formatting
	cachedResponseBytes, _ := json.Marshal(expectedProperty)
	cachedResponse := string(cachedResponseBytes)

	// Mock Redis Get behavior (simulate cache hit)
	mockRedis.On("Get", mock.Anything, "property123").Return(redis.NewStringResult(cachedResponse, nil))

	// Call the function
	result, err := api.GetPropertyDetail(mockRedis, nil, "property123", "user123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedProperty, result)

	// Verify expectations
	mockRedis.AssertExpectations(t)
}

func TestGetPropertyDetailAPIError(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockHTTP := new(MockHTTPClient)

	// Mock Redis Get behavior (simulate cache miss)
	mockRedis.On("Get", mock.Anything, "property123").Return(redis.NewStringResult("", redis.Nil))

	// Mock HTTP error response
	mockResp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Bad Request"}`)),
	}
	mockHTTP.On("Do", mock.Anything).Return(mockResp, nil)

	// Replace default HTTP client with mock
	originalHTTPClient := api.DefaultHTTPClient
	api.DefaultHTTPClient = mockHTTP
	defer func() { api.DefaultHTTPClient = originalHTTPClient }()

	// Call the function
	_, err := api.GetPropertyDetail(mockRedis, nil, "property123", "user123")

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API returned non-200 status code")

	// Verify expectations
	mockRedis.AssertExpectations(t)
	mockHTTP.AssertExpectations(t)
}
