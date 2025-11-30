package jwt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/DangeL187/erax"
	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func (m *Manager) Generate(id uint, tokenType string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": id,
		"typ": tokenType,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(m.privateKey)
}

func (m *Manager) GetPublicKey() (string, error) {
	pubDER, err := x509.MarshalPKIXPublicKey(m.publicKey)
	if err != nil {
		return "", erax.Wrap(err, "failed to marshal public key")
	}

	return base64.StdEncoding.EncodeToString(pubDER), nil
}

func (m *Manager) ParseToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("invalid signing method")
		}
		return m.publicKey, nil
	})

	if err != nil || !token.Valid {
		return 0, "", erax.Wrap(err, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid claims format")
	}

	subFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, "", errors.New("sub claim is missing or not a number")
	}

	typ, ok := claims["typ"].(string)
	if !ok {
		return 0, "", errors.New("typ claim is missing or not a string")
	}

	return uint(subFloat), typ, nil
}

func NewJWTManager() (*Manager, error) {
	privateB64 := os.Getenv("JWT_PRIVATE_KEY")
	publicB64 := os.Getenv("JWT_PUBLIC_KEY")

	if privateB64 == "" || publicB64 == "" {
		return nil, errors.New("missing JWT_PRIVATE_KEY or JWT_PUBLIC_KEY")
	}

	privateDER, err := base64.StdEncoding.DecodeString(privateB64)
	if err != nil {
		return nil, err
	}
	publicDER, err := base64.StdEncoding.DecodeString(publicB64)
	if err != nil {
		return nil, err
	}

	privateKeyIfc, err := x509.ParsePKCS8PrivateKey(privateDER)
	if err != nil {
		return nil, err
	}
	privateKey, ok := privateKeyIfc.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("not an ed25519 private key")
	}

	publicKeyIfc, err := x509.ParsePKIXPublicKey(publicDER)
	if err != nil {
		return nil, err
	}
	publicKey, ok := publicKeyIfc.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("not an ed25519 public key")
	}

	return &Manager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}
