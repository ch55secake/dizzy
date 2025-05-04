// Package output provides utility to properly output information to user
package output

import (
	"github.com/fatih/color"
)

func DefaultMessage() {
	asciiArt := `                ___
           ____/ (_)_______  __  __
          / __  / /_  /_  / / / / /
         / /_/ / / / /_/ /_/ /_/ /
         \__,_/_/ /___/___/\__, /
                          /____/
          An unsung hero.    `
	PrintCyanMessage(asciiArt, false)
	PrintCyanMessage("=================================================================", true)
}

func PrintCyanMessage(message string, square bool) {
	cyan := color.New(color.FgCyan, color.Bold)
	if square {
		_, err := cyan.Printf("[+] " + message + "\n")
		if err != nil {
			return
		}
	} else {
		_, err := cyan.Printf(message + "\n")
		if err != nil {
			return
		}
	}
}

func PrintMagentaMessage(message string, square bool) {
	magenta := color.New(color.FgMagenta, color.Bold)
	if square {
		_, err := magenta.Printf("[+] " + message + "\n")
		if err != nil {
			return
		}
	} else {
		_, err := magenta.Printf(message + "\n")
		if err != nil {
			return
		}
	}
}
