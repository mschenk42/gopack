package gopack

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type ColorFormat struct {
}

const colorFormat = "\033[%d;%dm%s\033[0m"

func (c ColorFormat) BlueBold(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 1, 34, s)
	}
	return s
}

func (c ColorFormat) Blue(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 0, 34, s)
	}
	return s
}

func (c ColorFormat) CyanBold(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 1, 36, s)
	}
	return s
}

func (c ColorFormat) Cyan(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 0, 36, s)
	}
	return s
}

func (c ColorFormat) GreenBold(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 1, 32, s)
	}
	return s
}

func (c ColorFormat) Green(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 0, 32, s)
	}
	return s
}

func (c ColorFormat) RedBold(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 1, 31, s)
	}
	return s
}

func (c ColorFormat) Red(s string) string {
	if c.IsTerminal() {
		return fmt.Sprintf(colorFormat, 0, 31, s)
	}
	return s
}

func (c ColorFormat) IsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
