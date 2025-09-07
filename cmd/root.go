package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Context struct {
	Debug bool
}

var CLI struct {
	Debug bool `help:"Enable debug mode"`

	ApiKey ApiKeyCmd `cmd:"api-key" help:"manage api keys"`
}

func main() {
	fmt.Println("share-docs' CLI")

	ctx := kong.Parse(&CLI)
	err := ctx.Run(&Context{Debug: CLI.Debug})

	ctx.FatalIfErrorf(err)
}
