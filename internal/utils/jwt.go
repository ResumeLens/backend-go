package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(viper.GetString("JWT_SECRET"))
}

type JWTClaim struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token
func GenerateJWT(userID, email, role string) (string, error) {
	expiryMinutes := viper.GetInt("JWT_EXPIRY")
	claims := &JWTClaim{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiryMinutes))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken parses and validates a token string
func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
