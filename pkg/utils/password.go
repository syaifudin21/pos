package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with its plaintext equivalent.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashOTP hashes the given OTP using bcrypt.
func HashOTP(otp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckOTPHash compares a hashed OTP with its plaintext equivalent.
func CheckOTPHash(otp, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(otp))
	return err == nil
}
