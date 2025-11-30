package token

import (
	"time"
)

type Generator interface {
	Generate(id uint, tokenType string, ttl time.Duration) (string, error)
}

type Manager interface {
	Generate(id uint, tokenType string, ttl time.Duration) (string, error)
	ParseToken(tokenString string) (uint, string, error)
}
