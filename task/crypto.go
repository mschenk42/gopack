package task

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mschenk42/gopack/task/internal/cryptopasta"
)

// Encrypt encrypts data using 256-bit AES-GCM. If no key is provided it is created.
// The encryption key and encrypted value are encoded as either hex or base64 and returned to the caller.
func Encrypt(text, encryptionKey string, base64Encoded bool) (string, string, error) {
	var key = &[32]byte{}

	if encryptionKey == "" {
		key = cryptopasta.NewEncryptionKey()
	} else {
		var err error
		key, err = decodeEncryptionKey(encryptionKey, base64Encoded)
		if err != nil {
			return "", "", fmt.Errorf("error decoding encryptiong key, %s", err)
		}
	}

	encrypted, err := cryptopasta.Encrypt([]byte(text), key)
	if err != nil {
		return "", "", fmt.Errorf("error encrypting text, %s", err)
	}

	return encodeValue(encrypted, base64Encoded), encodeEncryptionKey(key, base64Encoded), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.
func Decrypt(encrypted, encryptionKey string, base64Encoded bool) (string, error) {
	if encryptionKey == "" {
		return "", errors.New("encryption key is blank")
	}

	key, err := decodeEncryptionKey(encryptionKey, base64Encoded)
	if err != nil {
		return "", fmt.Errorf("error decoding encryptiong key, %s", err)
	}

	b, err := decodeValue(encrypted, base64Encoded)
	if err != nil {
		return "", fmt.Errorf("error decoding value, %s", err)
	}

	decrypted, err := cryptopasta.Decrypt(b, key)
	if err != nil {
		return "", fmt.Errorf("error decrypting value, %s", err)
	}

	return string(decrypted), nil
}

func encodeEncryptionKey(b *[32]byte, base64Key bool) string {
	if base64Key {
		return base64.StdEncoding.EncodeToString(b[:])
	}
	return hex.EncodeToString(b[:])
}

func encodeValue(b []byte, base64Value bool) string {
	if base64Value {
		return base64.StdEncoding.EncodeToString(b[:])
	}
	return hex.EncodeToString(b[:])
}

func decodeEncryptionKey(s string, base64Key bool) (*[32]byte, error) {
	k := [32]byte{}
	b := []byte{}
	var err error
	if base64Key {
		b, err = base64.StdEncoding.DecodeString(s)
	} else {
		b, err = hex.DecodeString(s)
	}
	if err != nil {
		return &k, err
	}
	if len(b) != 32 {
		return &k, fmt.Errorf("encryption key size is not 32")
	}
	copy(k[:], b)
	return &k, nil
}

func decodeValue(s string, base64Value bool) ([]byte, error) {
	if base64Value {
		return base64.StdEncoding.DecodeString(s)
	}
	return hex.DecodeString(s)
}
