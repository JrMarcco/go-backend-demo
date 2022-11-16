package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPasswd(passwd string) (string, error) {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPasswd), nil
}

func CheckPasswd(passwd, hashedPasswd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd))
}
