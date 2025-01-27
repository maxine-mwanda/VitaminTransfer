package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/plutov/paypal/v4"
)

var db *sql.DB

type PaymentResponse struct {
	Status  string
	Message string
}

// InitializePayPalClient sets up the PayPal client
func InitializePayPalClient() (*paypal.Client, error) {
	clientID := os.Getenv("PAYPAL_CLIENT_ID")
	clientSecret := os.Getenv("PAYPAL_CLIENT_SECRET")
	isLive := os.Getenv("PAYPAL_ENV")

	// Create PayPal client
	client, err := paypal.NewClient(clientID, clientSecret, isLive)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PayPal client: %w", err)
	}

	// Get access token
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get PayPal access token: %w", err)
	}

	return client, nil
}

func savePaymentToDB(order *paypal.Order) error {
	query := `
        INSERT INTO payments (id, intent, payer_email, total, currency, description) 
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(
		query,
		order.ID,
		order.Intent,
		order.Payer.EmailAddress,
		order.PurchaseUnits[0].Amount.Value,
		order.PurchaseUnits[0].Amount.Currency,
		order.PurchaseUnits[0].Description,
	)
	if err != nil {
		return fmt.Errorf("failed to save payment to database: %w", err)
	}
	return nil
}

func ProcessPayPalPayment(amount float64, currency string) (*PaymentResponse, error) {
	client, err := InitializePayPalClient()
	if err != nil {
		return nil, err
	}

	// Create a PayPal order
	order := paypal.CreateOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []paypal.PurchaseUnitRequest{{
			Amount: &paypal.PurchaseUnitAmount{
				Value:    fmt.Sprintf("%.2f", amount),
				Currency: currency,
			},
			Description: "Donation to Vitamin Transfer",
		}},
		ApplicationContext: &paypal.ApplicationContext{
			ReturnURL: "http://localhost:8080/success",
			CancelURL: "http://localhost:8080/cancel",
		},
	}

	createdOrder, err := client.CreateOrder(context.Background(), order)
	if err != nil {
		return nil, fmt.Errorf("PayPal order creation failed: %w", err)
	}

	// Extract the approval URL for the user to approve the order
	var approvalURL string
	for _, link := range createdOrder.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	// Save payment to the database
	if err := savePaymentToDB(createdOrder); err != nil {
		return nil, err
	}

	return &PaymentResponse{
		Status:  "Success",
		Message: fmt.Sprintf("Please complete the payment by visiting: %s", approvalURL),
	}, nil
}

func saveVisaTransactionToDB(cardNumber, expiryDate string, amount float64, currency string) error {
	query := `
        INSERT INTO visa_transactions (card_number, expiry_date, amount, currency) 
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, cardNumber, expiryDate, amount, currency)
	if err != nil {
		return fmt.Errorf("failed to save Visa transaction to database: %w", err)
	}
	return nil
}

func ProcessVisaPayment(cardNumber, expiryDate, cvv string, amount float64, currency string) (*PaymentResponse, error) {
	// Normally, you would send these details to Visa's payment gateway
	// For security, card details should never be logged or stored
	visaGatewayURL := os.Getenv("VISA_API_URL")
	apiKey := os.Getenv("VISA_API_KEY")

	requestPayload := map[string]interface{}{
		"card_number": cardNumber,
		"expiry_date": expiryDate,
		"cvv":         cvv,
		"amount":      amount,
		"currency":    currency,
	}

	if _, err := json.Marshal(requestPayload); err != nil {
		return nil, fmt.Errorf("failed to serialize Visa payment payload: %w", err)
	}

	req, err := http.NewRequest("POST", visaGatewayURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create Visa payment request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Body = http.NoBody

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Visa payment request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Visa payment gateway responded with status: %s", resp.Status)
	}
	// Save successful transaction to the database
	if err := saveVisaTransactionToDB(cardNumber, expiryDate, amount, currency); err != nil {
		return nil, err
	}

	return &PaymentResponse{
		Status:  "Success",
		Message: "Visa payment processed successfully",
	}, nil
}

func ProcessMpesaPayment(amount string, phoneNumber string) error {
	consumerKey := os.Getenv("MPESA_CONSUMER_KEY")
	consumerSecret := os.Getenv("MPESA_CONSUMER_SECRET")
	if consumerKey == "" || consumerSecret == "" {
		return fmt.Errorf("M-Pesa credentials are missing")
	}
	// Placeholder: Integrate M-Pesa API logic
	fmt.Printf("Processing M-Pesa payment of %s for phone %s...\n", amount, phoneNumber)
	return nil
}
