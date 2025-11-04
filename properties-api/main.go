package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Properties API - Coming soon...")
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"properties-api"}`))
	})
	
	fmt.Println("Server starting on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
