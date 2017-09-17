package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/task"
)

func encrypt(kv, encryptionKey, propertyFile string, base64Key bool) error {
	parts := strings.Split(kv, "=")
	if len(parts) != 2 {
		return fmt.Errorf("key value pair %s not valid", kv)
	}
	key := parts[0]
	value := parts[1]

	if base64Key {
		var err error
		encryptionKey, err = task.Base64ToHex(encryptionKey)
		if err != nil {
			return fmt.Errorf("unable to convert encryption key for base64 to hex, %s", err)
		}
	}

	if encryptionKey == "" {
		encryptionKey = task.EncodeEncryptionKey(task.NewEncryptionKey())
	}

	encrypted, err := task.EncryptText(value, encryptionKey)
	if err != nil {
		return err
	}

	if propertyFile != "" {

		if !strings.HasSuffix(propertyFile, ".json") {
			return fmt.Errorf("unable to update property file %s, file does not have extension json", propertyFile)
		}

		p := &gopack.Properties{}
		_, exists, err := task.Fexists(propertyFile)
		if exists {
			f1, err := os.Open(propertyFile)
			if err != nil {
				return err
			}
			defer f1.Close()

			b, err := ioutil.ReadAll(f1)
			if err != nil {
				return err
			}

			// load properties into map
			// err = p.Load(f1)
			err = p.Load(bytes.NewBuffer(b))
			if err != nil {
				return err
			}

			// let's make a backup, since we will be rewriting the property file
			f2, err := os.Create(propertyFile + ".bak")
			if err != nil {
				return err
			}
			defer f2.Close()

			x, err := io.Copy(f2, bytes.NewBuffer(b))
			fmt.Println("number bytes written", x)
			if err != nil {
				return err
			}
		}

		// create new property file
		f, err := os.Create(propertyFile)
		if err != nil {
			return err
		}
		defer f.Close()

		// set password and save property file
		(*p)[key] = hex.EncodeToString(encrypted)
		err = p.Save(f)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n%s: %s\nencryption key: %s\n", key, hex.EncodeToString(encrypted), encryptionKey)
	return nil
}

func decrypt(encrypted, encryptionKey, propertyFile string, base64Key bool) error {
	if base64Key {
		var err error
		encryptionKey, err = task.Base64ToHex(encryptionKey)
		if err != nil {
			return fmt.Errorf("unable to convert encryption key for base64 to hex, %s", err)
		}

	}
	v, err := task.DecryptText(encrypted, encryptionKey)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s\n", v)
	return nil
}
