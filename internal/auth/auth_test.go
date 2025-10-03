package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSimplePassword(t *testing.T) {
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

func TestCorrectJWT(t *testing.T) {
	userId := uuid.New()
	jwt, err := MakeJWT(userId, "secret1", 10*time.Second)
	if err != nil {
		t.Errorf("unable to create jwt. err: %s", err)
	}

	jwtUserId, err := ValidateJWT(jwt, "secret1")
	if err != nil {
		t.Errorf("unable to check password. err: %s", err)
	}

	if jwtUserId != userId {
		t.Errorf("userId did not match. was %s; expected %s", jwtUserId, userId)
	}
}

func TestJWTDifferentSecret(t *testing.T) {
	userId := uuid.New()
	jwt, err := MakeJWT(userId, "secret1", 10*time.Second)
	if err != nil {
		t.Errorf("unable to create jwt. err: %s", err)
	}

	jwtUserId, err := ValidateJWT(jwt, "secret2")
	if err == nil {
		t.Errorf("expected token not to be valid because secret is different")
	}

	if jwtUserId == userId {
		t.Errorf("userId was returned on invalid token")
	}
}

func TestExpiredJWT(t *testing.T) {
	userId := uuid.New()
	jwt, err := MakeJWT(userId, "secret1", -1*time.Second)
	if err != nil {
		t.Errorf("unable to create jwt. err: %s", err)
	}

	jwtUserId, err := ValidateJWT(jwt, "secret1")
	if err == nil {
		t.Errorf("jwt was valid when expected invalid")
	}

	if jwtUserId == userId {
		t.Errorf("userId was returned on invalid token")
	}
}

func TestGetBearerToken(t *testing.T) {
	correctValue := "secretToken"
	headers := http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", correctValue)},
	}
	res, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("unable to get token from headers")
	}

	if res != correctValue {
		t.Errorf("header was incorrect. was: %s; expected %s", res, correctValue)
	}
}

func TestGetBearerToken_InvalidHeader(t *testing.T) {
	correctValue := "secretToken"
	headers := http.Header{
		"Authorization": {fmt.Sprintf("Bearwithme %s", correctValue)},
	}
	_, err := GetBearerToken(headers)
	if err == nil {
		t.Errorf("header was invalid but received no err")
	}
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	headers := http.Header{}
	_, err := GetBearerToken(headers)
	if err == nil {
		t.Errorf("header was invalid but received no err")
	}
}
