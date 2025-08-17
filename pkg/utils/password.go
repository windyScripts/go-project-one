package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(password, encodedHash string) error {
	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		ErrorHandler(errors.New("invalid encoded hash format"), "internal server error")
	}

	saltBase64 := parts[0]
	hashBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)

	if err != nil {
		return ErrorHandler(err, "internal error")
		// http.Error(w, "Failed to decode salt", http.StatusInternalServerError)
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(hashBase64)

	if err != nil {
		return ErrorHandler(err, "internal error")
		// http.Error(w, "Failed to decode the hashed password.", http.StatusInternalServerError)
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if len(hash) != len(hashedPassword) {
		return ErrorHandler(errors.New("hash length mismatch"), "incorrect password")
		// http.Error(w, "Incorrect password.", http.StatusForbidden)
	}

	if subtle.ConstantTimeCompare(hash, hashedPassword) != 1 {
		return ErrorHandler(errors.New("incorrect username or password"), "incorrect password")
		// http.Error(w, "Incorrect password.", http.StatusForbidden)
	}
	return nil
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is blank")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}
