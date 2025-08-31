package auth

import (
	"testing"
)

func TestPasswordPolicy_ValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		config   PasswordPolicyConfig
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid strong password",
			config: PasswordPolicyConfig{
				MinLength:    8,
				MaxLength:    72,
				MinUpper:     1,
				MinLower:     1,
				MinDigit:     1,
				MinSpecial:   1,
				SpecialChars: "!@#$%^&*()_+-=[]{}|;:,.<>?",
			},
			password: "StrongPass123!",
			wantErr:  false,
		},
		{
			name: "password too short",
			config: PasswordPolicyConfig{
				MinLength: 8,
			},
			password: "short",
			wantErr:  true,
			errMsg:   "password is too short",
		},
		{
			name: "password too long",
			config: PasswordPolicyConfig{
				MinLength: 8,
				MaxLength: 10,
			},
			password: "verylongpassword",
			wantErr:  true,
			errMsg:   "password is too long",
		},
		{
			name: "insufficient uppercase",
			config: PasswordPolicyConfig{
				MinLength: 8,
				MaxLength: 50,
				MinUpper:  2, // Require 2 uppercase letters
			},
			password: "lowercase123A", // Only 1 uppercase
			wantErr:  true,
			errMsg:   "password does not have enough uppercase letters",
		},
		{
			name: "insufficient lowercase",
			config: PasswordPolicyConfig{
				MinLength: 8,
				MaxLength: 50,
				MinLower:  3, // Require 3 lowercase letters
			},
			password: "UPPERCASEab123", // Only 2 lowercase
			wantErr:  true,
			errMsg:   "password does not have enough lowercase letters",
		},
		{
			name: "insufficient digits",
			config: PasswordPolicyConfig{
				MinLength: 8,
				MaxLength: 50,
				MinDigit:  2, // Require 2 digits
			},
			password: "NoDigitsHere1", // Only 1 digit
			wantErr:  true,
			errMsg:   "password does not have enough digits",
		},
		{
			name: "insufficient special characters",
			config: PasswordPolicyConfig{
				MinLength:    8,
				MaxLength:    50,
				MinSpecial:   2, // Require 2 special characters
				SpecialChars: "!@#$",
			},
			password: "NoSpecialChars123!", // Only 1 special character
			wantErr:  true,
			errMsg:   "password does not have enough special characters",
		},
		{
			name: "valid password with multiple requirements",
			config: PasswordPolicyConfig{
				MinLength:    12,
				MaxLength:    50,
				MinUpper:     2,
				MinLower:     3,
				MinDigit:     2,
				MinSpecial:   1,
				SpecialChars: "!@#$",
			},
			password: "MyValidPass123!", // 2 upper, 9 lower, 3 digits, 1 special
			wantErr:  false,
		},
		{
			name: "no requirements - any password valid",
			config: PasswordPolicyConfig{
				MinLength:  4,
				MaxLength:  100,
				MinUpper:   0,
				MinLower:   0,
				MinDigit:   0,
				MinSpecial: 0,
			},
			password: "simple",
			wantErr:  false,
		},
		{
			name: "complex requirements met",
			config: PasswordPolicyConfig{
				MinLength:    16,
				MaxLength:    50,
				MinUpper:     3,
				MinLower:     4,
				MinDigit:     3,
				MinSpecial:   2,
				SpecialChars: "!@#$%^&*",
			},
			password: "MyVerySecurePassword123!@", // 3 upper, 15 lower, 3 digits, 2 special
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewPasswordPolicyService(tt.config)
			err := service.ValidatePassword(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errMsg != "" && !containsSubstring(err.Error(), tt.errMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(str, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(str) < len(substr) {
		return false
	}

	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if str[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
