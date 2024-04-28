package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Hash(input string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), s.hashCost)
	if err != nil {
		return "", fmt.Errorf("cannot generate hash : %w", err)
	}

	return string(bytes), nil
}

func (s *service) CompareHash(input string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(input)) == nil
}
