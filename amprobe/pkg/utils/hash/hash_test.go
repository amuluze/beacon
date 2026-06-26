package hash

import (
	"testing"
)

func TestBcryptHashAndVerify(t *testing.T) {
	password := "my-secret-password"
	hashed, err := BcryptHash(password)
	if err != nil {
		t.Fatalf("BcryptHash failed: %v", err)
	}
	if hashed == "" {
		t.Fatal("BcryptHash returned empty string")
	}
	if hashed == password {
		t.Fatal("BcryptHash returned plaintext")
	}

	if err := BcryptVerify(password, hashed); err != nil {
		t.Fatalf("BcryptVerify should succeed for correct password: %v", err)
	}

	if err := BcryptVerify("wrong-password", hashed); err == nil {
		t.Fatal("BcryptVerify should fail for wrong password")
	}
}

func TestBcryptHashDifferentHashes(t *testing.T) {
	password := "same-password"
	h1, err := BcryptHash(password)
	if err != nil {
		t.Fatalf("BcryptHash failed: %v", err)
	}
	h2, err := BcryptHash(password)
	if err != nil {
		t.Fatalf("BcryptHash failed: %v", err)
	}
	if h1 == h2 {
		t.Fatal("bcrypt should produce different hashes for the same password")
	}
}

func TestConstantTimeCompare(t *testing.T) {
	if !ConstantTimeCompare("abc", "abc") {
		t.Fatal("ConstantTimeCompare should return true for equal strings")
	}
	if ConstantTimeCompare("abc", "abd") {
		t.Fatal("ConstantTimeCompare should return false for different strings")
	}
}
