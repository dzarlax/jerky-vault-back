package utils

import (
    "golang.org/x/crypto/bcrypt"
    "log"
)

// HashPassword хеширует пароль
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Ошибка при хешировании пароля:", err)
        return "", err
    }
    return string(bytes), nil
}

// CheckPassword проверяет соответствие пароля и хеша
func CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}