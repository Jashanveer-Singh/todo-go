package jwttoken

import (
	"time"

	"github.com/Jashanveer-Singh/todo-go/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewJWTTokenProvider(
	secretKey, issuer, audience string,
	validtiyPeriod time.Duration,
) jwtToken {
	return jwtToken{
		secretKey:      secretKey,
		issuer:         issuer,
		audience:       audience,
		validtiyPeriod: validtiyPeriod,
	}
}

type jwtToken struct {
	secretKey      string
	issuer         string
	audience       string
	validtiyPeriod time.Duration
}

func (jt jwtToken) GenerateToken(claims models.Claims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"iss":  jt.issuer,
		"aud":  jt.audience,
		"exp":  time.Now().Add(jt.validtiyPeriod).Unix(),
		"iat":  time.Now().Unix(),
		"id":   claims.ID,
		"role": claims.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(jt.secretKey))
}

func (jt jwtToken) ValidateToken(tokenString string) (models.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(jt.secretKey), nil
	})
	if err != nil {
		return models.Claims{}, err
	}

	if !token.Valid {
		return models.Claims{}, jwt.ErrTokenMalformed
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.Claims{}, jwt.ErrTokenInvalidClaims
	}

	if claims["iss"] != jt.issuer || claims["aud"] != jt.audience {
		return models.Claims{}, jwt.ErrTokenInvalidClaims
	}

	return jt.extractClaims(claims)
}

func (jt jwtToken) extractClaims(claims jwt.MapClaims) (models.Claims, error) {
	id, ok := claims["id"].(float64)
	if !ok {
		return models.Claims{}, jwt.ErrTokenInvalidClaims
	}

	role, ok := claims["role"].(string)
	if !ok {
		return models.Claims{}, jwt.ErrTokenInvalidClaims
	}

	return models.Claims{
		ID:   int64(id),
		Role: role,
	}, nil
}
