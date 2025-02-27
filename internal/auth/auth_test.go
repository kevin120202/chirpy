package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secret"
	hash, err := HashPassword(password)
	
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	
	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}
	
	if hash == password {
		t.Fatal("Hash should not be equal to the original password")
	}
	
	// Hash should be different each time due to salt
	secondHash, _ := HashPassword(password)
	if hash == secondHash {
		t.Fatal("Two hashes of the same password should be different due to salt")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "secret"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	
	// Correct password should verify
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash should not return error for correct password: %v", err)
	}
	
	// Incorrect password should fail verification
	wrongPassword := "wrong-password"
	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Fatal("CheckPasswordHash should return error for incorrect password")
	}
}