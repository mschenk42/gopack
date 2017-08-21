package main

import (
	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/example/mypack"
)

func main() {
	mypack.Run(gopack.ParseCLI())
}
