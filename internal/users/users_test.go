package users

import "testing"

func TestPasswordHashing(t *testing.T) {
	password := "testPassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hashedPassword) {
		t.Fatal("Password and hash do not match!")
	}

	if CheckPasswordHash("wrongPassword", hashedPassword) {
		t.Fatal("Wrong password matched the hash!")
	}
}
