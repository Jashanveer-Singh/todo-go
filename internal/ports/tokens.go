package ports

import "github.com/Jashanveer-Singh/todo-go/internal/models"

type TokenProvider interface {
	GenerateToken(claims models.Claims) (token string, err error)
	ValidateToken(token string) (claims models.Claims, err error)
}
