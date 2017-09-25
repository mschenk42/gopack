package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/task"
)

func encrypt(kv, encryptionKey, propertyFile string, base64Encoded bool) error {
	parts := strings.Split(kv, "=")
	if len(parts) != 2 {
		return fmt.Errorf("key value pair %s not valid", kv)
	}
	key := parts[0]
	value := parts[1]

	encrypted, encryptionKey, err := task.Encrypt(value, encryptionKey, base64Encoded)
	if err != nil {
		return err
	}

	if propertyFile != "" {

		if !strings.HasSuffix(propertyFile, ".json") {
			return fmt.Errorf("%s does not have json extension", propertyFile)
		}

		p := &gopack.Properties{}
		_, exists, err := task.Fexists(propertyFile)
		if exists {
			f1, err := os.Open(propertyFile)
			if err != nil {
				return err
			}
			defer func() {
				if err := f1.Close(); err != nil {
					panic(err)
				}
			}()

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
			defer func() {
				if err := f2.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f2, bytes.NewBuffer(b))
			if err != nil {
				return err
			}
		}

		// create new property file
		f, err := os.Create(propertyFile)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		// set password and save property file
		(*p)[key] = encrypted
		err = p.Save(f)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n%s: %s\nencryption key: %s\n", key, encrypted, encryptionKey)
	return nil
}
