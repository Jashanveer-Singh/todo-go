package bcrypt

import "golang.org/x/crypto/bcrypt"

func NewBcryptPasswordHasher(cost int) *bycrptPasswordHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.MinCost
	} else if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return &bycrptPasswordHasher{
		cost: cost,
	}
}

type bycrptPasswordHasher struct {
	cost int
}

func (b bycrptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b bycrptPasswordHasher) CompareHash(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
