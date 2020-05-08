package flags

import "github.com/urfave/cli/v2"

func GetConfigurationFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "path",
			Value:    "./configs",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "format",
			Value:    "yaml",
			Required: false,
		},
	}
}
