package ports

type PasswordHasher interface {
	Hash(password string) (hash string, err error)
	CompareHash(hash string, password string) (IsValid bool, err error)
}
