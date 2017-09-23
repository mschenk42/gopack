package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mschenk42/gopack/color"
)

const generateUsage = `
generate gopack packs and task
usage: gopack generate --type [pack|task] --name <string> file-path
`
const encryptUsage = `
encrypt property values with AES-256 using hex or base64 encoding
usage: gopack encrypt --key <string> --file <json file> --base64 key=value
`
const decryptUsage = `
decrypt property values with AES-256 using hex or base64 encoding
usage: gopack decrypt --key <string> --base64 encypted-value
`

var (
	help = flag.Bool("help", false, "general help")

	generateCommand = flag.NewFlagSet("generate", flag.ExitOnError)
	helpGenerate    = generateCommand.Bool("help", false, "generate command help")
	typeToGenerate  = generateCommand.String("type", "task", "task or pack")
	typeName        = generateCommand.String("name", "", "name of generated task or pack(defaults to path's base dir or file name)")

	encryptCommand      = flag.NewFlagSet("encrypt", flag.ExitOnError)
	helpEncrypt         = encryptCommand.Bool("help", false, "encrypt command help")
	propertyFileEncrypt = encryptCommand.String("file", "", "property file to update/add encrypted value")
	encryptKeyEncrypt   = encryptCommand.String("key", "", "key to use for encryption, defaults to hexadecimal encoded")
	base64KeyEncrypt    = encryptCommand.Bool("base64", false, "key is base64 encoded")

	decryptCommand    = flag.NewFlagSet("decrypt", flag.ExitOnError)
	helpDecrypt       = decryptCommand.Bool("help", false, "decrypt command help")
	decryptKeyDecrypt = decryptCommand.String("key", "", "key to use for decryption, defaults to hexadecimal encoded")
	base64KeyDecrypt  = decryptCommand.Bool("base64", false, "key is base64 encoded")
)

func main() {
	if *help {
		onError(nil)
	}

	if len(os.Args) < 2 {
		onError(errors.New("no subcommand provided"))
	}

	switch os.Args[1] {
	case "generate":
		if *helpGenerate {
			onErrorGenerate(nil)
		}
		if len(os.Args) < 3 {
			onErrorGenerate(fmt.Errorf("file path not provided"))
		}
		generateCommand.Parse(os.Args[2:])
		switch *typeToGenerate {
		case "task":
			if err := generateTask(*typeName, generateCommand.Arg(0), false); err != nil {
				onErrorGenerate(err)
			}
		case "pack":
			if err := generatePack(*typeName, generateCommand.Arg(0), false); err != nil {
				onErrorGenerate(err)
			}
		default:
			onErrorGenerate(fmt.Errorf("%s not valid", *typeToGenerate))
		}

	case "encrypt":
		if *helpEncrypt {
			onErrorEncrypt(nil)
		}
		if len(os.Args) < 3 {
			onErrorEncrypt(fmt.Errorf("key=value to encrypt not provided"))
		}
		encryptCommand.Parse(os.Args[2:])
		if err := encrypt(encryptCommand.Arg(0), *encryptKeyEncrypt, *propertyFileEncrypt, *base64KeyEncrypt); err != nil {
			onErrorEncrypt(err)
		}

	case "decrypt":
		if *helpDecrypt {
			onErrorDecrypt(nil)
		}
		if len(os.Args) < 3 {
			onErrorDecrypt(fmt.Errorf("key to unencrypt not provided"))
		}
		decryptCommand.Parse(os.Args[2:])
		if err := decrypt(decryptCommand.Arg(0), *decryptKeyDecrypt, *base64KeyDecrypt); err != nil {
			onErrorDecrypt(err)
		}

	default:
		onError(nil)
	}
}

func onError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, color.Red("%s"), err)
		fmt.Println()
	}
	fmt.Fprint(os.Stderr, generateUsage)
	fmt.Println()
	generateCommand.PrintDefaults()
	fmt.Fprint(os.Stderr, encryptUsage)
	fmt.Println()
	encryptCommand.PrintDefaults()
	fmt.Fprint(os.Stderr, decryptUsage)
	fmt.Println()
	decryptCommand.PrintDefaults()
	os.Exit(1)
}

func onErrorGenerate(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, color.Red("%s"), err)
		fmt.Println()
	}
	fmt.Fprint(os.Stderr, generateUsage)
	fmt.Println()
	generateCommand.PrintDefaults()
	os.Exit(1)
}

func onErrorEncrypt(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, color.Red("%s"), err)
		fmt.Println()
	}
	fmt.Fprint(os.Stderr, encryptUsage)
	fmt.Println()
	encryptCommand.PrintDefaults()
	os.Exit(1)
}

func onErrorDecrypt(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, color.Red("%s"), err)
		fmt.Println()
	}
	fmt.Fprint(os.Stderr, decryptUsage)
	fmt.Println()
	decryptCommand.PrintDefaults()
	os.Exit(1)
}

func confirm(prompt string) bool {
	response := ""
	fmt.Print(prompt)
	fmt.Scanln(&response)
	return strings.TrimSpace(strings.ToLower(response)) == "y"
}
