package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Search API - Coming soon...")
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"search-api"}`))
	})
	
	fmt.Println("Server starting on port 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
