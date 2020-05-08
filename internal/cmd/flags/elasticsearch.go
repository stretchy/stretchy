package flags

import "github.com/urfave/cli/v2"

func GetElasticSearchFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "elasticsearch-host",
			EnvVars:  []string{"ELASTICSEARCH_HOST"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "elasticsearch-user",
			EnvVars:  []string{"ELASTICSEARCH_USER"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "elasticsearch-password",
			EnvVars:  []string{"ELASTICSEARCH_PASSWORD"},
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "elasticsearch-debug",
			EnvVars:  []string{"ELASTICSEARCH_DEBUG"},
			Required: false,
		},
	}
}
