package main

import (
	"fmt"

	"github.com/mschenk42/gopack/task"
)

func decrypt(encrypted, encryptionKey string, base64Encoded bool) error {
	v, err := task.Decrypt(encrypted, encryptionKey, base64Encoded)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s\n", v)
	return nil
}
