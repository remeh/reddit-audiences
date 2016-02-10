package app

import (
	"golang.org/x/crypto/bcrypt"
)

// IsPasswordSecure checks that the given password
// is strong enough to be used.
func IsPasswordSecure(password string) bool {
	// TODO(remy): check the password force
	return true
}

func CryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(b), err
}
