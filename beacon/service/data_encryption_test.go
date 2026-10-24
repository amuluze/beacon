package service

import (
	"encoding/base64"
	"errors"
	"testing"
)

func TestNewDataCipherRequiresValidEnvironmentKey(t *testing.T) {
	t.Setenv(dataEncryptionKeyEnv, "")
	if _, err := NewDataCipher(&Config{}); !errors.Is(err, ErrDataEncryptionKeyRequired) {
		t.Fatalf("missing key error = %v, want ErrDataEncryptionKeyRequired", err)
	}
	t.Setenv(dataEncryptionKeyEnv, base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef")))
	if _, err := NewDataCipher(&Config{}); err != nil {
		t.Fatalf("valid key rejected: %v", err)
	}
}
