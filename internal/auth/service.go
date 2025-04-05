package auth

import (
	"file-sharing-system/pkg/jwtutil"
	"golang.org/x/crypto/bcrypt"
)
//  I have used for hashing password before saving
func HashPassword(password string) (string, error) {
	// 14 is cost value, can reduce if needed
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
// compares user input password with hashed one
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email string) (string, error) {
	return jwtutil.CreateToken(email)
}
