package password

import passVal "github.com/wagslane/go-password-validator"

const (
	minEntropy = 60
)

func ValidatePassword(password string) (bool, error) {
	err := passVal.Validate(password, minEntropy)
	if err != nil {
		return false, err
	}
	return true, nil
}
