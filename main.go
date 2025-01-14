package main

import (
	"log"
	"net/http"
	"VitaminTransfer/utils"
	"VitaminTransfer/controllers"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	fmt.Println("loaded .env")

	// Initialize logging
	utils.InitLogger()
	utils.Logger.Info("Starting Vitamin Transfer application...")

	// Define routes
	http.HandleFunc("/", controllers.HomeHandler)
	http.HandleFunc("/donate", controllers.DonateHandler)
	http.HandleFunc("/success", controllers.SuccessHandler)

	// Start server
	port := "8000"
	utils.Logger.Infof("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		utils.Logger.Fatalf("Error starting server: %v", err)
	}
}
