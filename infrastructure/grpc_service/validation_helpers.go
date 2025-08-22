package grpcservice

import (
	"fmt"
	"regexp"
)

// isValidEmail validates email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePasswordStrength validates password strength
func validatePasswordStrength(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("mật khẩu phải có ít nhất 6 ký tự")
	}

	if len(password) > 50 {
		return fmt.Errorf("mật khẩu không được vượt quá 50 ký tự")
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("mật khẩu phải chứa ít nhất 1 chữ hoa")
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("mật khẩu phải chứa ít nhất 1 chữ thường")
	}

	// Check for at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("mật khẩu phải chứa ít nhất 1 chữ số")
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		return fmt.Errorf("mật khẩu phải chứa ít nhất 1 ký tự đặc biệt")
	}

	return nil
}

// validatePasswordMatch validates if password and confirm password match
func validatePasswordMatch(password, confirmPassword string) error {
	if password != confirmPassword {
		return fmt.Errorf("mật khẩu và xác nhận mật khẩu không khớp")
	}
	return nil
}
