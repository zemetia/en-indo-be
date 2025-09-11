package helpers

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPassword(hashPassword string, plainPassword []byte) (bool, error) {
	hashPW := []byte(hashPassword)
	if err := bcrypt.CompareHashAndPassword(hashPW, plainPassword); err != nil {
		return false, err
	}
	return true, nil
}

// GeneratePasswordFromBirthDate generates password from birth date in format DDMMYYYY
// Example: 2003-11-15 → "15112003"
func GeneratePasswordFromBirthDate(birthDate time.Time) string {
	if birthDate.IsZero() {
		// Fallback to default password if birth date is invalid
		return "12345678"
	}

	// Format: DDMMYYYY
	return fmt.Sprintf("%02d%02d%d", birthDate.Day(), int(birthDate.Month()), birthDate.Year())
}
