package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const pepper = "qwf5yh8kza!&tg430KmaE4326bfkirs1"

func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Encrypt(text string, password string) (string, error) {
	key := []byte(password + pepper)[:32]

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("encrypt", err)
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encryptedText, password string) (string, error) {
	if encryptedText == "" {
		return "", errors.New("encryptedText is empty")
	}
	combinedPassword := []byte(password + pepper)[:32]
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "ciphertext error", err
	}
	block, err := aes.NewCipher(combinedPassword)
	if err != nil {
		return "block, err", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "gcm, err", err
	}
	nonSize := gcm.NonceSize()
	nonce, cipciphertext := ciphertext[:nonSize], ciphertext[nonSize:]

	plaintext, err := gcm.Open(nil, nonce, cipciphertext, nil)
	if err != nil {
		return "plaintext, err", err
	}
	return string(plaintext), nil

}
