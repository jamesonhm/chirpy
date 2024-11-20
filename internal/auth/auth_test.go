package auth

import (
	"testing"
)

func TestCheckHashed(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatal(err)
	}
}
