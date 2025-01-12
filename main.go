package main

import (
	"log"
	"net/http"
	"vitamin-transfer/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize logging
	utils.InitLogger()
	utils.Logger.Info("Starting Vitamin Transfer application...")

	// Define routes
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/donate", DonateHandler)
	http.HandleFunc("/success", SuccessHandler)

	// Start server
	port := "8000"
	utils.Logger.Infof("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		utils.Logger.Fatalf("Error starting server: %v", err)
	}
}
