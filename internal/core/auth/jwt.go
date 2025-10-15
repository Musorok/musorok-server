package auth

import (
	"time"
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func HashToken(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func NewAccessToken(secret string, userID, role string, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID, Role: role,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func ParseClaims(secret, token string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil { return nil, err }
	if c, ok := tok.Claims.(*Claims); ok && tok.Valid { return c, nil }
	return nil, jwt.ErrTokenInvalidClaims
}
