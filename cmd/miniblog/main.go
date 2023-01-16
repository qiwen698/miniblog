package main

import (
	"os"
)
import (
	"github.com/qiwen698/miniblog/internal/miniblog"
	_ "go.uber.org/automaxprocs"
)

func main() {
	command := miniblog.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
