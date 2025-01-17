package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
    return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err
}

func MakeJWT(
    userId uuid.UUID, 
    tokenSecret string, 
    expiresIn time.Duration,
) (string, error)  {
    now := time.Now().UTC()
    claims := jwt.RegisteredClaims{
        Issuer: "chirpy",
        IssuedAt: jwt.NewNumericDate(now),
        ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
        Subject: userId.String(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // overview of the different signing methods and their respective key types:
    // https://golang-jwt.github.io/jwt/usage/signing_methods/#signing-methods-and-key-types
    return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    claimsStruct := jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(
        tokenString, 
        &claimsStruct,
        func(token *jwt.Token) (interface{}, error) {
	        return []byte(tokenSecret), nil
        },
    )
    if err != nil {
        return uuid.Nil, err
    } 

    userIDString, err := token.Claims.GetSubject()
    if err != nil {
        return uuid.Nil, err
    } 
    
    id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

    return id, err
} 

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// GetBearerToken -
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// func MakeRefreshToken() (string, error) {
//     c := 32
//     b := make([]byte, c)
//     _, err := rand.Read(b)
//     if err != nil {
//         return "", err
//     }
//     return hex.EncodeToString(b), nil
// }

// MakeRefreshToken makes a random 256 bit token
// encoded in hex
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
