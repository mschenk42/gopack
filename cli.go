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
		fmt.Fprint(os.Stdout, "Run pack (y/n)? ")
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

		for idx, a := range args {
			fmt.Fprintf(os.Stdout, "loading %s configuration file\n", a)
			p2 := Properties{}

			f, err := os.Open(a)
			if err != nil {
				exitOnError(fmt.Errorf("unable to load property file %s, %s\n", a, err))
			}
			defer f.Close()

			if err := p2.Load(f); err != nil {
				exitOnError(fmt.Errorf("unable to load property file %s, %s\n", a, err))
			}
			if idx > 0 {
				p.Merge(&p2)
			} else {
				p = p2
			}
		}
		return &p
	}

	return nil
}

func exitOnError(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}
