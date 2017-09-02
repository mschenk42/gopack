package gopack

import (
	"fmt"
	"os"

	//TODO: how much does this increase the exe size?
	"golang.org/x/crypto/ssh/terminal"
)

type ColorFormat struct {
}

func (c ColorFormat) GreenBold(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf("\033[%d;%dm%s\033[0m", 1, 32, s)
	}
	return s
}

func (c ColorFormat) Green(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf("\033[%d;%dm%s\033[0m", 0, 32, s)
	}
	return s
}

func (c ColorFormat) RedBold(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf("\033[%d;%dm%s\033[0m", 1, 31, s)
	}
	return s
}

func (c ColorFormat) Red(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf("\033[%d;%dm%s\033[0m", 0, 31, s)
	}
	return s
}

func (c ColorFormat) IsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
