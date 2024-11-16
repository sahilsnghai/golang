package auth

import "testing"

func TestJWTToken(t *testing.T) {

	hash, err := HashPassword("passowrd")
	if err != nil {
		t.Errorf("found error while hashin : %s", err)
	}

	if hash == "" {
		t.Error("Hash expected not to be empty string")

	}

	if hash == "password" {
		t.Error("Hash expected not to be same as original password")

	}
}
