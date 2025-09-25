package main

import (
	"fmt"
	"os"

	toolcmd "github.com/NSObjects/go-template/tools/cmd"
)

func main() {
	if err := toolcmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
