package user

import passwordvalidator "github.com/wagslane/go-password-validator"

type PasswordEntropyResult struct {
	Entropy    float64
	MinEntropy float64
	Valid      bool
}

func (s *service) PasswordEntropy(password string) PasswordEntropyResult {
	minEntropy := float64(s.passwordMinEntropy.Get())
	entropy := passwordvalidator.GetEntropy(password)
	return PasswordEntropyResult{
		Entropy:    entropy,
		MinEntropy: minEntropy,
		Valid:      entropy >= minEntropy,
	}
}
