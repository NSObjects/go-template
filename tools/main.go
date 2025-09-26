package main

import (
	"fmt"
	"os"

	toolcmd "github.com/NSObjects/go-template/tools/cmd"
)

func main() {
	// 设置命令名
	cmd := toolcmd.NewRootCommand()
	cmd.Use = "gt"

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
