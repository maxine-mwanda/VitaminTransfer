package main

import (
	"VitaminTransfer/controllers"
	"VitaminTransfer/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using system environment variables.")
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
	port := os.Getenv("PORT")
	utils.Logger.Infof("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		utils.Logger.Fatalf("Error starting server: %v", err)
	}
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))
}
