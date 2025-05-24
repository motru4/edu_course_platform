package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	pepper string
}

func NewHasher(pepper string) *Hasher {
	return &Hasher{
		pepper: pepper,
	}
}

func (h *Hasher) Hash(password string) (string, error) {
	peppered := h.addPepper(password)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(peppered), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hashedBytes), nil
}

func (h *Hasher) Compare(password, hash string) error {
	peppered := h.addPepper(password)

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(peppered))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}

func (h *Hasher) addPepper(password string) string {
	return password + h.pepper
}
