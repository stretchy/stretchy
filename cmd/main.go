package main

import (
	"log"
	"os"

	"github.com/stretchy/stretchy/internal/cmd/apply"
	"github.com/urfave/cli/v2"
)

var version = "development"

func main() {
	app := &cli.App{
		Name:    "stretchy",
		Usage:   "Elasticsearch migrations",
		Version: version,
		Commands: []*cli.Command{
			apply.GetApplyCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
