package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"net/http"

	"github.com/plutov/paypal/v4" // PayPal SDK
)


type PaymentResponse struct {
	Status  string
	Message string
}

// InitializePayPalClient sets up the PayPal client
func InitializePayPalClient() (*paypal.Client, error) {
	clientID := os.Getenv("PAYPAL_CLIENT_ID")
	clientSecret := os.Getenv("PAYPAL_CLIENT_SECRET")
	isLive := os.Getenv("PAYPAL_ENV") == "live"

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

func ProcessPayPalPayment(amount float64, currency string) (*PaymentResponse, error) {
	client, err := InitializePayPalClient()
	if err != nil {
		return nil, err
	}

	// Create a PayPal payment
	payment := paypal.Payment{
		Intent: "sale",
		Payer: &paypal.Payer{
			PaymentMethod: "paypal",
		},
		Transactions: []paypal.Transaction{{
			Amount: &paypal.Amount{
				Total:    fmt.Sprintf("%.2f", amount),
				Currency: currency,
			},
			Description: "Donation to Vitamin Transfer",
		}},
		RedirectURLs: &paypal.RedirectURLs{
			ReturnURL: "http://localhost:8080/success",
			CancelURL: "http://localhost:8080/cancel",
		},
	}

	createdPayment, err := client.CreatePayment(context.Background(), payment)
	if err != nil {
		return nil, fmt.Errorf("PayPal payment creation failed: %w", err)
	}

	// Extract the approval URL for the user to approve the payment
	var approvalURL string
	for _, link := range createdPayment.Links {
		if link.Rel == "approval_url" {
			approvalURL = link.Href
			break
		}
	}

	return &PaymentResponse{
		Status:  "Pending",
		Message: fmt.Sprintf("Please complete the payment by visiting: %s", approvalURL),
	}, nil
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

	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
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