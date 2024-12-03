package fns

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"io"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/hooks/slack/errors"
)

// encrypt encrypts the given plaintext using the provided key and AES-256 in Cipher Feedback mode.
// It returns the encrypted ciphertext, or an error if encryption fails.
func Encrypt(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	return ciphertext, nil
}

// decrypt decrypts the given ciphertext using the provided key and AES-256 in Cipher Feedback mode.
// It returns the decrypted plaintext, or an error if decryption fails.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.ErrCipherText
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// generate creates a 32-byte key for AES-256 encryption.
// It uses a SHA-512 hash of the provided workspaceID,
// and then takes the first 32 bytes of the hash as the key.
func Generate(workspaceID string) []byte {
	h := sha512.New()
	h.Write([]byte(auth.Secret() + workspaceID)) // TODO - verify

	return h.Sum(nil)[:32]
}