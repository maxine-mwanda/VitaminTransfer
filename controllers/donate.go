package controllers

import (
	"VitaminTransfer/models"
	"VitaminTransfer/utils"
	"fmt"
	"net/http"
)

func DonateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	amount := r.FormValue("amount")
	paymentMethod := r.FormValue("paymentMethod")

	utils.Logger.Infof("Processing donation: Name=%s, Email=%s, Amount=%s, PaymentMethod=%s", name, email, amount, paymentMethod)

	var err error

	switch paymentMethod {
	case "PayPal":
		err = models.ProcessPayPalPayment(amount, "USD")
	case "Visa":
		cardNumber := r.FormValue("cardNumber")
		expiry := r.FormValue("expiry")
		cvv := r.FormValue("cvv")
		err = models.ProcessVisaPayment(amount, cardNumber, expiry, cvv)
	case "Mpesa":
		phoneNumber := r.FormValue("phoneNumber")
		err = models.ProcessMpesaPayment(amount, phoneNumber)
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
