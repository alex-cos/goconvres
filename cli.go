package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

var (
	version   = "unknown"
	buildDate = "unknown"
)

// parseCLI parses command lines arguments.
func parseCLI() error {
	cliapp := cli.NewApp()
	cliapp.Name = "goconvres"
	cliapp.Usage = "Command line to convert resource file (PNG,JPEG,Binary,...) into go source file"
	cliapp.UsageText = "goconvres [options] <input file> <output file>"
	cliapp.Description = "Build: " + buildDate
	cliapp.Version = version

	cliapp.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Printf("Error. Unknown command: '%s'\n\n", command)
		cli.ShowAppHelpAndExit(c, 1)
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("Version:\t", c.App.Version)
	}

	cliapp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "package",
			Value:       "resources",
			Usage:       "Package name for resources",
			Aliases:     []string{"p"},
			Required:    false,
			DefaultText: fmt.Sprintf("Default value is '%s'", "resources"),
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "Resource",
			Usage:       "Specify the name of the resource",
			Aliases:     []string{"n"},
			Required:    false,
			DefaultText: fmt.Sprintf("Default value is '%s'", "Resource"),
		},
		&cli.IntFlag{
			Name:        "ncols",
			Value:       12,
			Usage:       "Number of columns to format",
			Aliases:     []string{"c"},
			Required:    false,
			DefaultText: fmt.Sprintf("Default value is '%d'", 12),
		},
	}

	cliapp.Action = action

	sort.Sort(cli.FlagsByName(cliapp.Flags))
	sort.Sort(cli.CommandsByName(cliapp.Commands))

	if err := cliapp.Run(os.Args); err != nil {
		return fmt.Errorf("failed to parse command line arguments: %w", err)
	}
	return nil
}
