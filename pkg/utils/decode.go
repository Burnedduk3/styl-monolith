package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"styl-monolith/internal/users/core/domain"
)

func DecodeJWT(token string) domain.JwtClaims {
	// Split the token into its parts (header, payload, signature)
	parts := strings.Split(token, ".")
	if len(parts) < 3 {
		fmt.Println("Invalid JWT token")
		return domain.JwtClaims{}
	}

	// Decode Header (optional)
	_, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		fmt.Println("Error decoding header:", err)
		return domain.JwtClaims{}
	}
	//fmt.Println("Header:", string(header))

	// Decode Payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("Error decoding payload:", err)
		return domain.JwtClaims{}
	}
	//fmt.Println("Payload:", string(payload))

	// Unmarshal payload JSON to Go struct (optional)
	var claims domain.JwtClaims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		fmt.Println("Error unmarshaling claims:", err)
		return domain.JwtClaims{}
	}
	return claims
}
