package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"` // Add role here
	jwt.StandardClaims
}

// GenerateJWT creates a token with an additional role claim.
func GenerateJWT(email string, username string, role string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1000 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		Role:     role, // Set role here
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

// ValidateToken validates the token and checks expiration.
func ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	// Remove Bearer prefix if present
	if len(signedToken) > 7 && signedToken[:7] == "Bearer " {
		signedToken = signedToken[7:]
	}

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	parsedClaims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if parsedClaims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token expired")
	}
	return parsedClaims, nil
}
