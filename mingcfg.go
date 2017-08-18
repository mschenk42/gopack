package mincfg

import (
	"log"
	"os"
)

func init() {
	LogInfo = log.New(os.Stdout, "INFO: ", 0)
	LogErr = log.New(os.Stderr, "ERROR: ", 0)
}

var (
	LogInfo *log.Logger
	LogErr  *log.Logger
)
