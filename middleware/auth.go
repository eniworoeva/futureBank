package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(tokenString, "Bearer ")

		_, err := VerifyJWT(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})

}

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header is required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	parsedtoken, err := VerifyJWT(tokenString)
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := parsedtoken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("could not get claims")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return 0, fmt.Errorf("userID not found in claims")
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return 0, err
	}

	return uint(userIDInt), nil

}
