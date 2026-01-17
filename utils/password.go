package utils

import (
    "golang.org/x/crypto/bcrypt"
    "log"
)

// HashPassword hashes the password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        return "", err
    }
    return string(bytes), nil
}

// CheckPassword verifies the password matches the hash
func CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}