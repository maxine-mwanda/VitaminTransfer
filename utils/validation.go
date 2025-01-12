package utils

import "regexp"

func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func IsValidAmount(amount string) bool {
	amountRegex := `^[0-9]+(\.[0-9]{1,2})?$`
	return regexp.MustCompile(amountRegex).MatchString(amount)
}