package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestJWT(t *testing.T) {
	const testDelay time.Duration = time.Second * 2
	goodSecret := "longstringconvertedtobyteslater"
	badSecret := "toosimple"
	user1 := uuid.New()
	user2 := uuid.New()

	tests := []struct {
		name        string
		userIn      uuid.UUID
		secretIn    string
		userValid   uuid.UUID
		secretValid string
		expiresIn   time.Duration
		wantErr     bool
		match       bool
	}{
		{
			name:        "Correct All",
			userIn:      user1,
			secretIn:    goodSecret,
			userValid:   user1,
			secretValid: goodSecret,
			expiresIn:   time.Second * 5,
			wantErr:     false,
			match:       true,
		},
		{
			name:        "Incorrect secret",
			userIn:      user1,
			secretIn:    goodSecret,
			userValid:   user1,
			secretValid: badSecret,
			expiresIn:   time.Second * 5,
			wantErr:     true,
			match:       false,
		},
		{
			name:        "Incorrect user",
			userIn:      user1,
			secretIn:    goodSecret,
			userValid:   user2,
			secretValid: goodSecret,
			expiresIn:   time.Second * 5,
			wantErr:     false,
			match:       false,
		},
		{
			name:        "time expire",
			userIn:      user1,
			secretIn:    goodSecret,
			userValid:   user2,
			secretValid: goodSecret,
			expiresIn:   time.Second * 1,
			wantErr:     true,
			match:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userIn, tt.secretIn, tt.expiresIn)
			if err != nil {
				t.Errorf("error making JWT: %v", err)
			}
			time.Sleep(testDelay)
			id, err := ValidateJWT(token, tt.secretValid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Failed to validate token, err = %v", err)
			}
			if (id == tt.userValid) != tt.match {
				t.Errorf("User id's do not match, in = %v, out = %v, expected = %v", tt.userIn, id, tt.userValid)
			}
		})
	}
}
