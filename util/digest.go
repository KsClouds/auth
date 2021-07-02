package util

import (
	"golang.org/x/crypto/bcrypt"
)

const PassWordCost = 12

func Digest(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	return string(bytes), err
}

func CheckPassword(password string, digest string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(digest), []byte(password))
	return err == nil
}
