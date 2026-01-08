package auth

import (
	
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    // Cost: 12-14 is common (higher = slower/more secure). DefaultCost is 10.
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
