package hashing

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type HashProvider struct {
	cost int
}

const _defaultHashLen = 32

func NewHashProvider(cost int) *HashProvider {
	return &HashProvider{cost: cost}
}

func (h *HashProvider) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (h *HashProvider) VerifyPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GenerateRandomToken generates a crypto-resistant random token.
// Length in bytes: for example, 32 bytes → 43 characters in Base64.
func (h *HashProvider) GenerateRandomToken() string {
	return h.generateRandomToken(_defaultHashLen)
}

// generateRandomToken with a given length (in bytes).
func (h *HashProvider) generateRandomToken(byteLen int) string {
	b := make([]byte, byteLen)

	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random bytes: " + err.Error())
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

// Hash creates the SHA-256 hash from the line and returns it in HEX format.
func (h *HashProvider) Hash(data string) string {
	hash := sha256.Sum256([]byte(data))

	return fmt.Sprintf("%x", hash)
}
