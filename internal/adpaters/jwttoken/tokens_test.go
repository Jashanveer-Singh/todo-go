package jwttoken

import (
	"testing"
	"time"

	"github.com/Jashanveer-Singh/todo-go/internal/models"
)

func Test_jwttoken_GenerateToken(t *testing.T) {
	jwtTokenProvider := NewJWTTokenProvider("mysecretkey", "myissuer", "myaudience", time.Hour)
	claims := models.Claims{
		ID:   1,
		Role: "",
	}

	token, err := jwtTokenProvider.GenerateToken(claims)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected token to be non-empty")
	}
}

func Test_jwttoken_ValidateToken_when_valid(t *testing.T) {
	jwtTokenProvider := NewJWTTokenProvider("mysecretkey", "myissuer", "myaudience", time.Hour)
	claims := models.Claims{
		ID:   1,
		Role: "",
	}

	token, err := jwtTokenProvider.GenerateToken(claims)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	validatedClaims, err := jwtTokenProvider.ValidateToken(token)
	if err != nil {
		t.Errorf("expected no error validating token, got %v", err)
	}
	if validatedClaims.ID != claims.ID {
		t.Errorf("expected UserID to be %d, got %d", claims.ID, validatedClaims.ID)
	}
	if validatedClaims.Role != claims.Role {
		t.Errorf("expected Role to be %s, got %s", claims.Role, validatedClaims.Role)
	}
}

func Test_jwttoken_ValidateToken_when_signed_with_invalid_secret(t *testing.T) {
	claims := models.Claims{
		ID:   1,
		Role: "",
	}

	jwtTokenProvider := NewJWTTokenProvider("mysecretkey", "myissuer", "myaudience", time.Hour)
	token, err := jwtTokenProvider.GenerateToken(claims)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	anotherJWTTokenProvider := NewJWTTokenProvider(
		"anothersecretkey",
		"myissuer",
		"myaudience",
		time.Hour,
	)

	_, err = anotherJWTTokenProvider.ValidateToken(token)
	if err == nil {
		t.Fatalf("expected error while validating token, got %v", err)
	}
}

func Test_jwttoken_ValidateToken_when_signed_with_invalid_iss_or_aud(t *testing.T) {
	claims := models.Claims{
		ID:   1,
		Role: "",
	}

	jwtTokenProvider := NewJWTTokenProvider("mysecretkey", "myissuer", "myaudience", time.Hour)
	token, err := jwtTokenProvider.GenerateToken(claims)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	anotherJWTTokenProvider := NewJWTTokenProvider(
		"mysecretkey",
		"anotherissuer",
		"myaudience",
		time.Hour,
	)

	_, err = anotherJWTTokenProvider.ValidateToken(token)
	if err == nil {
		t.Errorf("expected error while validating token, got %v", err)
	}

	anotherJWTTokenProvider = NewJWTTokenProvider(
		"mysecretkey",
		"myissuer",
		"anotheraudience",
		time.Hour,
	)
	_, err = anotherJWTTokenProvider.ValidateToken(token)
	if err == nil {
		t.Errorf("expected error while validating token, got %v", err)
	}
}
