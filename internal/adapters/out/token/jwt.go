package token

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/JoelTeoGom/login-service/internal/ports"
)

// JWT implements ports.TokenIssuer with HS256.
type JWT struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

var _ ports.TokenIssuer = (*JWT)(nil)

func NewJWT(secret string, ttl time.Duration, issuer string) *JWT {
	return &JWT{secret: []byte(secret), ttl: ttl, issuer: issuer}
}

func (j *JWT) Issue(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Issuer:    j.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *JWT) Verify(raw string) (int64, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(raw, claims,
		func(*jwt.Token) (any, error) { return j.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), // prevents the "alg: none" attack
		jwt.WithIssuer(j.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, errors.New("invalid subject")
	}
	return id, nil
}
