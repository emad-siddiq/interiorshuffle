package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// RedisClient interface defines the methods we need from redis
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient is the default client used for HTTP requests
var DefaultHTTPClient HTTPClient = &http.Client{}

// PropertyDetail represents the property detail response from the API
type PropertyDetail struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Price   string `json:"price"`
}

// GetPropertyDetail fetches property details from cache, DB, or API
func GetPropertyDetail(redisClient RedisClient, db *gorm.DB, propertyID string, userID string) (*PropertyDetail, error) {
	// Check the cache first
	cachedDetails, err := getPropertyDetailsFromCache(redisClient, propertyID)
	if err == nil && cachedDetails != nil {
		return cachedDetails, nil
	}

	// If not found in cache, check the DB (not implemented here but can be done similarly)
	// If not found in DB, fetch from the API
	details, err := callPropertyAPI(propertyID, userID, DefaultHTTPClient)
	if err != nil {
		return nil, fmt.Errorf("error calling API: %v", err)
	}

	// Cache the details
	err = cachePropertyDetails(redisClient, propertyID, details)
	if err != nil {
		return nil, fmt.Errorf("error caching property details: %v", err)
	}

	// Return the fetched details
	return details, nil
}

// callPropertyAPI makes an external API request to fetch property details
func callPropertyAPI(propertyID string, userID string, httpClient HTTPClient) (*PropertyDetail, error) {
	url := "https://api.realestateapi.com/v2/PropertyDetail"
	payload := map[string]interface{}{
		"id":    propertyID,
		"comps": false,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", "test")
	req.Header.Set("x-user-id", userID)

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	// Parse the response
	var propertyDetail PropertyDetail
	if err := json.NewDecoder(resp.Body).Decode(&propertyDetail); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &propertyDetail, nil
}

// getPropertyDetailsFromCache retrieves property details from Redis cache
func getPropertyDetailsFromCache(client RedisClient, propertyID string) (*PropertyDetail, error) {
	details, err := client.Get(context.Background(), propertyID).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("error fetching property details from cache: %v", err)
	}

	var propertyDetail PropertyDetail
	if err := json.Unmarshal([]byte(details), &propertyDetail); err != nil {
		return nil, fmt.Errorf("error unmarshaling cached details: %v", err)
	}

	return &propertyDetail, nil
}

// cachePropertyDetails stores property details in Redis cache for one hour
func cachePropertyDetails(client RedisClient, propertyID string, details *PropertyDetail) error {
	jsonData, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("error marshaling property details: %v", err)
	}

	err = client.Set(context.Background(), propertyID, string(jsonData), time.Hour).Err()
	if err != nil {
		return fmt.Errorf("could not cache property details: %v", err)
	}
	return nil
}
