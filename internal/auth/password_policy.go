package auth

import (
	"fmt"
	"regexp"
	"unicode"
)

var (
	ErrPasswordTooShort            = fmt.Errorf("password is too short")
	ErrPasswordTooLong             = fmt.Errorf("password is too long")
	ErrPasswordInsufficientUpper   = fmt.Errorf("password does not have enough uppercase letters")
	ErrPasswordInsufficientLower   = fmt.Errorf("password does not have enough lowercase letters")
	ErrPasswordInsufficientDigit   = fmt.Errorf("password does not have enough digits")
	ErrPasswordInsufficientSpecial = fmt.Errorf("password does not have enough special characters")
	ErrPasswordInvalidChar         = fmt.Errorf("password contains invalid characters")
)

// PasswordPolicyService defines the interface for password policy validation
type PasswordPolicyService interface {
	ValidatePassword(password string) error
}

type passwordPolicyService struct {
	config PasswordPolicyConfig
}

// NewPasswordPolicyService creates a new password policy service with the given configuration
func NewPasswordPolicyService(config PasswordPolicyConfig) PasswordPolicyService {
	return &passwordPolicyService{
		config: config,
	}
}

// ValidatePassword validates a password against the configured policy
func (p *passwordPolicyService) ValidatePassword(password string) error {
	// Check minimum length
	minLen := int(p.config.MinLength)
	if len(password) < minLen {
		return fmt.Errorf("%w: minimum length is %d characters", ErrPasswordTooShort, minLen)
	}

	// Check maximum length (0 means no limit)
	maxLen := int(p.config.MaxLength)
	if maxLen > 0 && len(password) > maxLen {
		return fmt.Errorf("%w: maximum length is %d characters", ErrPasswordTooLong, maxLen)
	}

	// Check for uppercase letters
	minUpper := int(p.config.MinUpper)
	if minUpper > 0 {
		upperCount := countUpper(password)
		if upperCount < minUpper {
			return fmt.Errorf("%w: requires %d, found %d", ErrPasswordInsufficientUpper, minUpper, upperCount)
		}
	}

	// Check for lowercase letters
	minLower := int(p.config.MinLower)
	if minLower > 0 {
		lowerCount := countLower(password)
		if lowerCount < minLower {
			return fmt.Errorf("%w: requires %d, found %d", ErrPasswordInsufficientLower, minLower, lowerCount)
		}
	}

	// Check for digits
	minDigit := int(p.config.MinDigit)
	if minDigit > 0 {
		digitCount := countDigit(password)
		if digitCount < minDigit {
			return fmt.Errorf("%w: requires %d, found %d", ErrPasswordInsufficientDigit, minDigit, digitCount)
		}
	}

	// Check for special characters
	minSpecial := int(p.config.MinSpecial)
	if minSpecial > 0 {
		specialChars := string(p.config.SpecialChars)
		specialCount := countSpecial(password, specialChars)
		if specialCount < minSpecial {
			return fmt.Errorf("%w: requires %d from %s, found %d", ErrPasswordInsufficientSpecial, minSpecial, specialChars, specialCount)
		}
	}

	return nil
}

// Helper functions for character counting

func countUpper(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsUpper(r) {
			count++
		}
	}
	return count
}

func countLower(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLower(r) {
			count++
		}
	}
	return count
}

func countDigit(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			count++
		}
	}
	return count
}

func countSpecial(s, specialChars string) int {
	if specialChars == "" {
		// If no special chars defined, use common special characters
		specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	count := 0
	// Create a regex pattern to match any of the special characters
	pattern := "[" + regexp.QuoteMeta(specialChars) + "]"
	re, _ := regexp.Compile(pattern)
	matches := re.FindAllString(s, -1)
	if matches != nil {
		count = len(matches)
	}
	return count
}
