package helpers

import "golang.org/x/crypto/bcrypt"

func HashPass(p string) string {
	password := []byte(p)

	hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	return string(hash)
}

func ComparePass(h, p []byte) bool {
	// hash, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(h, p)

	return err == nil
}
