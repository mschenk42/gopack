package gopack

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	yesFlag = flag.Bool("y", false, "run pack without confirmation")
)

const (
	defaultProperties = "gopack.json"
)

func LoadProperties() *Properties {
	flag.Parse()

	confirm := *yesFlag
	if !confirm {
		response := ""
		fmt.Print("Run pack (y/n)? ")
		fmt.Scanln(&response)
		confirm = strings.TrimSpace(strings.ToLower(response)) == "y"
	}

	if confirm {
		p := Properties{}

		args := flag.Args()
		if len(args) == 0 {
			if _, err := os.Stat(defaultProperties); err == nil {
				args = append(args, defaultProperties)
			}
		}

		for idx, f := range args {
			fmt.Printf("loading %s configuration file\n", f)
			p2 := Properties{}
			if err := p2.Load(f); err != nil {
				fmt.Fprintf(os.Stderr, "unable to load property file %s, %s\n", f, err)
				os.Exit(1)
				return nil
			}
			if idx > 0 {
				p.Merge(&p2)
			} else {
				p = p2
			}
		}
		return &p
	}

	os.Exit(0)
	return nil
}
