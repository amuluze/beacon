// Package hash
// Date: 2024/3/27 16:52
// Author: Amu
// Description:
package hash

import (
	"crypto/sha1"
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// SHA1String returns the SHA-1 hex digest of the input string.
func SHA1String(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

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
