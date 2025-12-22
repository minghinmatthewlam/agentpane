package main

import (
	"fmt"
	"os"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/cmd"
)

func main() {
	a, err := app.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := cmd.NewRootCmd(a).Execute(); err != nil {
		os.Exit(1)
	}
}
