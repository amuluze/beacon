package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

const encryptedPrefix = "enc:v1:"

var (
	ErrInvalidKey        = errors.New("data encryption key must decode to exactly 32 bytes")
	ErrInvalidCiphertext = errors.New("invalid encrypted secret")
)

type Cipher struct {
	aead cipher.AEAD
}

func NewFromEncodedKey(encoded string) (*Cipher, error) {
	key, err := decodeKey(strings.TrimSpace(encoded))
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKey
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{aead: aead}, nil
}

func decodeKey(encoded string) ([]byte, error) {
	if encoded == "" {
		return nil, ErrInvalidKey
	}
	decoders := []func(string) ([]byte, error){
		hex.DecodeString,
		base64.StdEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
	}
	for _, decode := range decoders {
		if key, err := decode(encoded); err == nil && len(key) == 32 {
			return key, nil
		}
	}
	return nil, ErrInvalidKey
}

func (c *Cipher) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := c.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptedPrefix + base64.RawStdEncoding.EncodeToString(sealed), nil
}

// Decrypt returns the plaintext and whether the input was legacy plaintext.
// Callers use the legacy flag to transparently rewrite old database rows.
func (c *Cipher) Decrypt(value string) (plaintext string, legacy bool, err error) {
	if value == "" {
		return "", false, nil
	}
	if !strings.HasPrefix(value, encryptedPrefix) {
		return value, true, nil
	}
	encoded := strings.TrimPrefix(value, encryptedPrefix)
	payload, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil || len(payload) < c.aead.NonceSize()+c.aead.Overhead() {
		return "", false, ErrInvalidCiphertext
	}
	nonce := payload[:c.aead.NonceSize()]
	ciphertext := payload[c.aead.NonceSize():]
	plaintextBytes, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", false, fmt.Errorf("%w", ErrInvalidCiphertext)
	}
	return string(plaintextBytes), false, nil
}

func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}
