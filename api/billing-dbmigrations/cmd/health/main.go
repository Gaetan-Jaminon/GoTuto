package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	
	_ "github.com/lib/pq"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database struct {
		Connected bool   `json:"connected"`
		Version   string `json:"version,omitempty"`
		Error     string `json:"error,omitempty"`
	} `json:"database"`
}

func main() {
	port := os.Getenv("HEALTH_PORT")
	if port == "" {
		port = "8081"
	}
	
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	
	log.Printf("Health check server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start health server:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "healthy"}
	
	// Check database connection
	db, err := getDB()
	if err != nil {
		response.Database.Connected = false
		response.Database.Error = err.Error()
	} else {
		defer db.Close()
		
		// Test connection
		if err := db.Ping(); err != nil {
			response.Database.Connected = false
			response.Database.Error = err.Error()
		} else {
			response.Database.Connected = true
			
			// Get migration version
			var version int
			err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
			if err == nil {
				response.Database.Version = fmt.Sprintf("%d", version)
			}
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// For readiness, we just check if we can connect to the database
	db, err := getDB()
	if err != nil {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func getDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "billing")
	sslmode := getEnv("DB_SSL_MODE", "disable")
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	
	return sql.Open("postgres", dsn)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}