package password

import "testing"

func TestPasswordHelper(t *testing.T) {
	weakPassword := "aweakpass"
	strongPasasword := "Str0nGP@ssword"
	res, err := ValidatePassword(weakPassword)
	if res {
		t.Errorf(err.Error())
	}
	res, err = ValidatePassword(strongPasasword)
	if !res {
		t.Errorf(err.Error())
	}
}
