package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		Text          string
		EncryptionKey string
		Base64encoded bool
		Error         bool
	}{
		{
			"test",
			"7c17a24a8124c47a5ed042950389ccc5bea46f66006a0a474bb37de965aacd00",
			false,
			false,
		},
		{
			"test",
			"aFR3AyKpbJYwFdyfax2RK2/62ctUoJHOkJC8oqtYJ/M=",
			true,
			false,
		},
		{
			"test",
			"",
			false,
			false,
		},
		{
			"test",
			"",
			true,
			false,
		},
		{
			`89^&23898!@#!SFAFdafdfljaflkjdfIIidafd||||||\\\}[]`,
			"",
			true,
			false,
		},
		{
			"test",
			"7c17a24a8124",
			false,
			true,
		},
	}

	for _, tc := range cases {
		var text string
		var err error
		encrypted, key, err := Encrypt(tc.Text, tc.EncryptionKey, tc.Base64encoded)
		if tc.Error {
			assert.EqualError(err, "error decoding encryptiong key, encryption key size is not 32")
		} else {
			assert.NoError(err, "error calling encrypt")
			if tc.EncryptionKey == "" {
				text, err = Decrypt(encrypted, key, tc.Base64encoded)
			} else {
				text, err = Decrypt(encrypted, tc.EncryptionKey, tc.Base64encoded)
			}
			assert.NoError(err, "error calling decrypt")
			assert.Equal(tc.Text, text, "decrypted text does not match")
		}
	}
}
