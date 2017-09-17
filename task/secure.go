// cryptopasta - basic cryptography examples
//
// Written in 2015 by George Tankersley <george.tankersley@gmail.com>
//
// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along
// with this software. If not, see // <http://creativecommons.org/publicdomain/zero/1.0/>.

// Provides symmetric authenticated encryption using 256-bit AES-GCM with a random nonce.
package task

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

func EncodeEncryptionKey(b *[32]byte) string {
	return hex.EncodeToString(b[:])
}

func DecodeEncryptionKey(s string) (*[32]byte, error) {
	k := [32]byte{}
	b, err := hex.DecodeString(s)
	if err != nil {
		return &k, err
	}
	copy(k[:], b)
	return &k, nil
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

func EncryptText(text, encryptionKey string) ([]byte, error) {
	var key = &[32]byte{}

	if encryptionKey == "" {
		key = NewEncryptionKey()
	} else {
		var err error
		key, err = DecodeEncryptionKey(encryptionKey)
		if err != nil {
			return []byte{}, fmt.Errorf("error decoding encryptiong key, %s", err)
		}
	}

	encrypted, err := Encrypt([]byte(text), key)
	if err != nil {
		return []byte{}, fmt.Errorf("error encrypting text, %s", err)
	}

	return encrypted, nil
}

func DecryptText(encrypted, encryptionKey string) (string, error) {
	if encryptionKey == "" {
		return "", errors.New("encryption key is blank")
	}

	byteKey, err := DecodeEncryptionKey(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("error decoding encryptiong key, %s", err)
	}

	b, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("error decoding value, %s", err)
	}

	decrypted, err := Decrypt(b, byteKey)
	if err != nil {
		return "", fmt.Errorf("error decrypting value, %s", err)
	}

	return string(decrypted), nil
}

func Base64ToHex(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
