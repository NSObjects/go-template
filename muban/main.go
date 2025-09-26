package main

import (
	"fmt"
	"os"

	"github.com/NSObjects/go-template/muban/cmd"
)

func main() {
	xdd := cmd.NewRootCommand()
	if err := xdd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
