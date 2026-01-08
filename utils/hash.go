package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPassword compares plain password with hashed password
func CheckPassword(hashed, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
    return err == nil
}
