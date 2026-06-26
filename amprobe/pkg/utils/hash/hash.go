// Package hash
// Date: 2024/3/27 16:52
// Author: Amu
// Description:
package hash

import (
	"crypto/subtle"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash generates a bcrypt hash from the given password.
func BcryptHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// BcryptVerify compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success, or an error on failure.
func BcryptVerify(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ConstantTimeCompare compares two strings in constant time to prevent timing attacks.
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
