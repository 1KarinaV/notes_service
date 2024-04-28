package auth

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"net/http"
	"notes_service/internal/services"
	"time"
)

func (s *service) MakeToken(email string) (string, error) {
	claims := map[string]interface{}{services.Login: email}
	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour*time.Duration(s.tokenLifetime)))

	_, tokenString, err := s.tkn.Encode(claims)
	if err != nil {
		return "", fmt.Errorf("cannot generate token : %w", err)
	}

	return tokenString, nil
}

func (s *service) GetClaims(request *http.Request) (string, error) {
	claims := jwtauth.TokenFromHeader(request)

	if claims == "" {
		return "", services.ErrEmptyClaims
	}

	token, err := s.tkn.Decode(claims)
	if err != nil {
		return "", fmt.Errorf("cannot decode token : %w", err)
	}

	if email, rs := token.Get(services.Login); rs {
		return email.(string), nil
	}

	return "", services.ErrEmptyClaims
}
