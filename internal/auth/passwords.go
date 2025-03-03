package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) error {
	bytesPassword := []byte(password)
	bytesHash := []byte(hash)

	err := bcrypt.CompareHashAndPassword(bytesHash, bytesPassword)
	if err != nil {
		return err
	}

	return nil
}
