package flags

import "github.com/urfave/cli/v2"

func MigrateOptions() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "enable-soft-update",
			Usage:   "Enable inplace remapping whenever it's possible",
			EnvVars: []string{"ENABLE_SOFT_UPDATE"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			EnvVars: []string{"DRY_RUN"},
			Value:   false,
		},
	}
}
