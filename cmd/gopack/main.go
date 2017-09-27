package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const generateUsage = `
usage: gopack generate --type <pack|task> --name <name> <file>
`
const encryptUsage = `
usage: gopack encrypt --key <key> --file <file> --base64 <key=value>
`
const decryptUsage = `
usage: gopack decrypt --key <key> --base64 <encypted string>
`

var (
	generateFlags  = flag.NewFlagSet("generate", flag.ContinueOnError)
	typeToGenerate = generateFlags.String("type", "task", "generate a task or pack")
	typeName       = generateFlags.String("name", "", "task or pack name, defaults to path's base dir or file name")

	encryptFlags     = flag.NewFlagSet("encrypt", flag.ContinueOnError)
	fileEncrypt      = encryptFlags.String("file", "", "property file to add encrypted value")
	keyEncrypt       = encryptFlags.String("key", "", "key to use for encryption")
	base64KeyEncrypt = encryptFlags.Bool("base64", false, "key is base64 encoded otherwise defaults to hex encoding")

	decryptFlags     = flag.NewFlagSet("decrypt", flag.ContinueOnError)
	keyDecrypt       = decryptFlags.String("key", "", "key to use for decryption")
	base64KeyDecrypt = decryptFlags.Bool("base64", false, "key is base64 encoded otherwise defaults to hex encoding")
)

func main() {
	generateFlags.Usage = func() {
		fmt.Fprintln(os.Stderr, generateUsage)
		generateFlags.PrintDefaults()
	}

	encryptFlags.Usage = func() {
		fmt.Fprintln(os.Stderr, encryptUsage)
		encryptFlags.PrintDefaults()
	}

	decryptFlags.Usage = func() {
		fmt.Fprintln(os.Stderr, decryptUsage)
		decryptFlags.PrintDefaults()
	}

	command := ""
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	switch command {
	case "generate":
		if err := generateFlags.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "file path not provided")
			generateFlags.Usage()
			os.Exit(1)
		}
		switch *typeToGenerate {
		case "task":
			if err := generateTask(*typeName, generateFlags.Arg(0), false); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		case "pack":
			if err := generatePack(*typeName, generateFlags.Arg(0), false); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		default:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s is not a valid type to generate\n", *typeToGenerate))
		}

	case "encrypt":
		if err := encryptFlags.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "key=value to encrypt not provided")
			encryptFlags.Usage()
			os.Exit(1)
		}
		if err := encrypt(encryptFlags.Arg(0), *keyEncrypt, *fileEncrypt, *base64KeyEncrypt); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "decrypt":
		if err := decryptFlags.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "key to decyrpt not provided")
			decryptFlags.Usage()
			os.Exit(1)
		}
		if err := decrypt(decryptFlags.Arg(0), *keyDecrypt, *base64KeyDecrypt); err != nil {
			os.Exit(1)
		}

	default:
		m := fmt.Sprintf("%s is not a valid command", command)
		if command == "" {
			m = "no command provided"
		}
		fmt.Fprintln(os.Stderr, m)
		os.Exit(1)
	}
}

func confirm(prompt string) bool {
	response := ""
	fmt.Print(prompt)
	fmt.Scanln(&response)
	return strings.TrimSpace(strings.ToLower(response)) == "y"
}
