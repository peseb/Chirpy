package auth

import "testing"

func TestSimpleInput(t *testing.T) {
	password := "Secret12!"
	hashed_password, err := HashPassword(password)
	if err != nil {
		t.Errorf("unable to hash password. err: %s", err)
	}

	password_correct, err := CheckPasswordHash(password, hashed_password)
	if err != nil {
		t.Errorf("unable to check password. err: %s", err)
	}

	if !password_correct {
		t.Errorf("password and hash did not match")
	}
}
