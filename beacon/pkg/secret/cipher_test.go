package secret

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

func testCipher(t *testing.T) *Cipher {
	t.Helper()
	key := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := NewFromEncodedKey(key)
	if err != nil {
		t.Fatalf("NewFromEncodedKey: %v", err)
	}
	return cipher
}

func TestCipherRoundTripAndLegacyDetection(t *testing.T) {
	cipher := testCipher(t)
	encrypted, err := cipher.Encrypt("sensitive-value")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !IsEncrypted(encrypted) || strings.Contains(encrypted, "sensitive-value") {
		t.Fatalf("encrypted value is not opaque: %q", encrypted)
	}
	plaintext, legacy, err := cipher.Decrypt(encrypted)
	if err != nil || legacy || plaintext != "sensitive-value" {
		t.Fatalf("Decrypt = %q, legacy=%v, err=%v", plaintext, legacy, err)
	}
	plaintext, legacy, err = cipher.Decrypt("old-plaintext")
	if err != nil || !legacy || plaintext != "old-plaintext" {
		t.Fatalf("legacy Decrypt = %q, legacy=%v, err=%v", plaintext, legacy, err)
	}
}

func TestCipherRejectsWrongKeyAndTampering(t *testing.T) {
	if _, err := NewFromEncodedKey("short"); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("key error = %v, want ErrInvalidKey", err)
	}
	cipher := testCipher(t)
	if _, _, err := cipher.Decrypt("enc:v1:not-base64"); !errors.Is(err, ErrInvalidCiphertext) {
		t.Fatalf("ciphertext error = %v, want ErrInvalidCiphertext", err)
	}
}
