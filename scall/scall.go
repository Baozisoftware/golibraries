//+build !windows

package scall

import "fmt"

func SetTerminalTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}
