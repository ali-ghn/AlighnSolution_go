package email

import "testing"

func TestEmailValidation(t *testing.T) {
	validEmail := "email@example.com"
	invalidEmail := "example.com"
	res, err := ValidateEmail(validEmail)
	if err != nil || !res {
		t.Errorf(err.Error())
	}

	res, err = ValidateEmail(invalidEmail)
	if err == nil || res {
		t.Errorf(err.Error())
	}

}
