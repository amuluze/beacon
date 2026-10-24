package auth

import (
	"errors"
	"unicode"
)

const MinPasswordLength = 12

var ErrWeakPassword = errors.New("password must be at least 12 characters and include at least three character classes")

// ValidatePassword enforces the password policy used for bootstrap accounts
// and password changes. Requiring three character classes avoids accepting a
// long but trivially guessable value while still allowing passphrases.
func ValidatePassword(password string) error {
	if len([]rune(password)) < MinPasswordLength {
		return ErrWeakPassword
	}

	classes := 0
	var lower, upper, digit, symbol bool
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			lower = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsDigit(r):
			digit = true
		default:
			symbol = true
		}
	}
	for _, present := range []bool{lower, upper, digit, symbol} {
		if present {
			classes++
		}
	}
	if classes < 3 {
		return ErrWeakPassword
	}
	return nil
}
