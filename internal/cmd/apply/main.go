package apply

import (
	"fmt"
	"path/filepath"

	"github.com/stretchy/stretchy/internal/cmd/flags"
	"github.com/stretchy/stretchy/pkg/action"
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/elasticsearch"
	"github.com/urfave/cli/v2"
)

func GetApplyCommand() *cli.Command {
	return &cli.Command{
		Name:  "apply",
		Usage: "",
		Flags: flags.Merge(
			flags.GetConfigurationFlags(),
			flags.GetElasticSearchFlags(),
			[]cli.Flag{
				&cli.StringSliceFlag{
					Name: "index-names",
				},
				&cli.StringFlag{
					Name:    "index-prefix",
					EnvVars: []string{"INDEX_PREFIX"},
				},
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
			},
		),
		Action: execute,
	}
}

func execute(c *cli.Context) error {
	indexCollection, err := load(c)
	if err != nil {
		return err
	}

	client, err := elasticsearch.New(
		elasticsearch.Options{
			Host:     c.String("elasticsearch-host"),
			User:     c.String("elasticsearch-user"),
			Password: c.String("elasticsearch-password"),
			Debug:    c.Bool("elasticsearch-debug"),
		},
	)

	if err != nil {
		return err
	}

	compareResultCollection, err := compare(
		c,
		indexCollection,
		client,
	)

	if err != nil {
		return err
	}

	fmt.Printf("Diffs:\n")

	for _, compareResult := range compareResultCollection {
		fmt.Printf("\tIndex '%s' => %s\n", compareResult.AliasName, compareResult.Result.Action().String())

		for _, d := range compareResult.Result.Changes() {
			fmt.Printf("\t\t%s\n", d.String())
		}
	}

	if c.Bool("dry-run") {
		return nil
	}

	return apply(client, compareResultCollection)
}

func load(c *cli.Context) (configuration.IndexCollection, error) {
	configPath, err := filepath.Abs(c.String("path"))
	if err != nil {
		return nil, err
	}

	loadAction := action.NewLoad(configPath)
	format := c.String("format")

	configurationNames := c.StringSlice("index-names")
	if len(configurationNames) == 0 {
		return loadAction.LoadAll(format)
	}

	indexCollection := configuration.IndexCollection{}

	for _, name := range configurationNames {
		index, err := loadAction.Load(name, c.String("format"))
		if err != nil {
			return nil, err
		}

		indexCollection.Load(name, index)
	}

	return indexCollection, nil
}

func compare(
	c *cli.Context,
	indexCollection configuration.IndexCollection,
	client elasticsearch.Client,
) (action.CompareResultCollection, error) {
	compareAction := action.NewCompare(client, c.String("index-prefix"), c.Bool("enable-soft-update"))

	return compareAction.CompareAll(indexCollection)
}

func apply(client elasticsearch.Client, compareResultCollection action.CompareResultCollection) error {
	applyAction := action.NewApply(client)

	return applyAction.ApplyAll(compareResultCollection)
}
