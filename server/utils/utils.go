package utils

import "regexp"

// validateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

// validatePassword checks if the password meets minimum strength requirements
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	re := regexp.MustCompile(`[a-zA-Z]`)
	return re.MatchString(password) && regexp.MustCompile(`[0-9]`).MatchString(password)
}


