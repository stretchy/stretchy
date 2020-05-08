package flags

import "github.com/urfave/cli/v2"

func Merge(flagGroups ...[]cli.Flag) []cli.Flag {
	flags := []cli.Flag{}

	for _, flagGroup := range flagGroups {
		flags = append(flags, flagGroup...)
	}

	return flags
}
