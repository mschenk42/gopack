package gopack

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultProperties = "gopack.json"
)

var (
	yesFlag    = flag.Bool("y", false, "run pack without confirmation")
	actionFlag = flag.String("actions", "", "pack actions to run")
	helpFlag   = flag.Bool("h", false, "show help and exit")
)

func usage() string {
	x := filepath.Base(os.Args[0])
	return fmt.Sprintf(`

%s [-y] [--actions action1,action2] [property1.json property2.json ...]

  * attempts to load "gopack.json" if no property files are specified
`, x)
}

func LoadProperties() (*Properties, []string) {
	flag.Parse()

	if *helpFlag {
		fmt.Fprint(os.Stdout, usage())
		flag.PrintDefaults()
		os.Exit(0)
	}

	actions := []string{}
	if strings.TrimSpace(*actionFlag) != "" {
		actions = strings.Split(*actionFlag, ",")
	}

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
			fmt.Fprintf(os.Stdout, PackPropertyFormat, fmt.Sprintf("loading %s configuration file", a))
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
		return &p, actions
	} else {
		os.Exit(0)
	}
	return nil, actions
}

func exitOnError(err error) {
	fmt.Fprint(os.Stderr, usage())
	flag.PrintDefaults()
	if err != nil {
		fmt.Fprintf(os.Stderr, PackErrorFormat, err)
	}
	os.Exit(1)
}
