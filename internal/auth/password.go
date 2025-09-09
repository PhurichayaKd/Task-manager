package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(pw string) (string, error)
	Compare(hash, pw string) error
}

type bcryptHasher struct{}

func NewPasswordHasher() PasswordHasher { return &bcryptHasher{} }

func (*bcryptHasher) Hash(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}
func (*bcryptHasher) Compare(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
