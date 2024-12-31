package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/go-redis/redis/v8"
	"interiorshuffle.com/property_app/api"
	"interiorshuffle.com/property_app/db"
)

// Initialize the Redis client
func InitRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
	return client
}

// PropertyDetail struct represents the property details.
type PropertyDetail struct {
	PropertyID string `json:"property_id"`
	Details    string `json:"details"`
}

func main() {
	//Load env variables
	err := godotenv.Load()

	// Initialize Redis client
	redisClient := InitRedisClient()
	// Initialize PostgreSQL database client
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}

	// Set up the HTTP handler for the property details endpoint
	http.HandleFunc("/property/details", func(w http.ResponseWriter, r *http.Request) {
		// Extract parameters from the query string
		propertyID := r.URL.Query().Get("property_id")
		userID := r.URL.Query().Get("user_id")
		if propertyID == "" || userID == "" {
			http.Error(w, "Missing property_id or user_id", http.StatusBadRequest)
			return
		}

		// Call the GetPropertyDetail function
		details, err := api.GetPropertyDetail(redisClient, db, propertyID, userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching property details: %v", err), http.StatusInternalServerError)
			return
		}

		// Marshal the details into JSON and write the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(details)
	})

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
