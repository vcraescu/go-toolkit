package cmd

import (
	"os"
)

func Terminate() {
	termCh <- os.Interrupt
}

func init() {
	termCh = make(chan os.Signal, 1)
}
