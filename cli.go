package gopack

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func ParseCLI() *Properties {
	yesFlag := flag.Bool("y", false, "run pack without confirmation")
	flag.Parse()

	confirm := *yesFlag
	if !confirm {
		response := ""
		fmt.Print("Run pack? ")
		fmt.Scanln(&response)
		confirm = strings.TrimSpace(strings.ToLower(response)) == "y"
	}

	if confirm {
		p := new(Properties)

		cfgFile := flag.Arg(0)
		if cfgFile == "" {
			cfgFile = "gopack.json"
		}

		fmt.Printf("loading %s configuration file\n", cfgFile)
		if err := p.Load(cfgFile); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load property file %s, %s", cfgFile, err)
			os.Exit(1)
			return nil
		}
		return p
	}

	os.Exit(0)
	return nil
}
