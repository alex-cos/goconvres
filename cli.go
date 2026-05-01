package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v3"
)

var (
	version   = "unknown"
	buildDate = "unknown"
)

// parseCLI parses command lines arguments.
func parseCLI() error {
	appcmd := &cli.Command{
		Name:        "goconvres",
		Usage:       "Command line to convert resource file (PNG,JPEG,Binary,...) into go source file",
		UsageText:   "goconvres [options] <input file> <output file>",
		Description: "Build: " + buildDate,
		Version:     version,
		CommandNotFound: func(c context.Context, cmd *cli.Command, name string) {
			fmt.Fprintf(os.Stderr, "Error. Unknown command: '%s'\n\n", name)
			cli.ShowAppHelpAndExit(cmd, 1)
		},
		Flags: []cli.Flag{
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
		},
	}

	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintln(os.Stdout, "Version:\t", cmd.Version)
	}

	appcmd.Action = action

	sort.Sort(cli.FlagsByName(appcmd.Flags))
	sort.Slice(appcmd.Commands, func(i, j int) bool {
		return appcmd.Commands[i].Name < appcmd.Commands[j].Name
	})

	if err := appcmd.Run(context.Background(), os.Args); err != nil {
		return fmt.Errorf("failed to parse command line arguments: %w", err)
	}

	return nil
}
