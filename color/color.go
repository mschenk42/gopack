package color

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const colorFormat = "\033[%d;%dm%s\033[0m"

func GreyBold(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 1, 37, s)
	}
	return s
}

func Grey(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 0, 37, s)
	}
	return s
}

func MagentaBold(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 1, 35, s)
	}
	return s
}

func Magenta(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 0, 35, s)
	}
	return s
}

func BlueBold(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 1, 34, s)
	}
	return s
}

func Blue(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 0, 34, s)
	}
	return s
}

func YellowBold(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 1, 33, s)
	}
	return s
}

func Yellow(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 0, 33, s)
	}
	return s
}

func CyanBold(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 1, 36, s)
	}
	return s
}

func Cyan(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 0, 36, s)
	}
	return s
}

func GreenBold(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 1, 32, s)
	}
	return s
}

func Green(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 0, 32, s)
	}
	return s
}

func RedBold(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 1, 31, s)
	}
	return s
}

func Red(s string) string {
	if isTerminal() {
		return fmt.Sprintf(colorFormat, 0, 31, s)
	}
	return s
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
