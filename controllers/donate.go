package controllers

import (
	"VitaminTransfer/models"
	"VitaminTransfer/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type DonationRequest struct {
	PaymentMethod string  `json:"paymentMethod"`
	Amount        float64 `json:"amount"`
	CardNumber    string  `json:"cardNumber"`
	Expiry        string  `json:"expiry"`
	CVV           string  `json:"cvv"`
	PhoneNumber   string  `json:"phoneNumber"`
}

var allowedMethods = map[string]bool{
	"PayPal": true,
	"Visa":   true,
	"Mpesa":  true,
}

func DonateHandler(w http.ResponseWriter, r *http.Request) {

var db *sql.DB

	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	// Read the raw body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		log.Println("Error reading request body:", err)
		return
	}

	log.Println("Raw request body:", string(body)) // Debugging log

	var donation DonationRequest

		// Now decode the JSON
	err = json.Unmarshal(body, &donation)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}
	response := map[string]string{"message": "Donation received successfully"}
	w.Header().Set("Content-Type", "application/json") // ✅ Ensure JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response) // ✅ Encode response correctly

	fmt.Println("Sending JSON response:", response)

	if err := json.NewEncoder(w).Encode(response); err != nil {
        fmt.Println("Error encoding response:", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }

	// Validate the payment method
	if !allowedMethods[donation.PaymentMethod] {
		return
	}
	//validate the data parsed
	if donation.PaymentMethod == "" || donation.Amount <= 0 {
		http.Error(w, "Invalid payment method or amount", http.StatusBadRequest)
		log.Printf("Received invalid data: Method=%s, Amount=%.2f\n", donation.PaymentMethod, donation.Amount)
		return
	}

	_, err = models.ProcessPayPalPayment(donation.Amount, "USD", db)

	utils.Logger.Infof("Processing donation: Amount=%s, PaymentMethod=%s", donation.Amount, donation.PaymentMethod)

	switch donation.PaymentMethod {
	case "PayPal":
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		_, err = models.ProcessPayPalPayment(donation.Amount, "USD", db)
	case "Visa":
		if donation.CardNumber == "" || donation.Expiry == "" || donation.CVV == "" {
			http.Error(w, "Missing Visa card details", http.StatusBadRequest)
			return
		}
		//_, err = models.ProcessVisaPayment(donation.Amount, donation.CardNumber, donation.Expiry, amountFloat, cvv, db *sql.DB)
	case "Mpesa":
		if donation.PhoneNumber == "" {
			http.Error(w, "Missing M-Pesa phone number", http.StatusBadRequest)
			return
		}
		//_,err = models.ProcessMpesaPayment(donation.Amount, donation.PhoneNumber)

	default:
		http.Error(w, "Invalid payment method", http.StatusBadRequest)
		return
	}

	if err != nil {
		utils.Logger.Errorf("Payment processing failed: %v", err)
		http.Error(w, fmt.Sprintf("Payment failed: %v", err), http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Payment processed successfully")

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}
