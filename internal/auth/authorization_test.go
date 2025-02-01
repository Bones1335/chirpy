package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	want, _ := bcrypt.GenerateFromPassword([]byte(password), 5)
	asString := string(want)
	hashedPass, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword('password') = %q, %v, want match for %#q, error", hashedPass, err, asString)
	}
}

