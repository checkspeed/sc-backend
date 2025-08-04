package utils

import (
	"regexp"
	"strings"
)

func IsValidText(text string, required bool) bool {
	text = strings.TrimSpace(text)
	if required && text == "" {
		return false
	}
	if text == "" { // Optional field
		return true
	}

	// Must contain at least one letter/number
	hasAlphaNum := regexp.MustCompile(`[a-zA-Z0-9]`).MatchString(text)
	// Not all special characters
	notAllSpecial := !regexp.MustCompile(`^[\s\W_]+$`).MatchString(text)

	return hasAlphaNum && notAllSpecial
}

func IsValidEmail(email string) bool {
	return regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`).MatchString(email)
}
