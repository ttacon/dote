package main

// DO NOT MODIFY: generated by github.com/ttacon/toml2cli
// If you need to make changes, update the toml config and then regenerate
// this file.

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "dote"
	app.Usage = "Manage and consume dotefiles"

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:        "list-profiles",
			Description: "List profiles",
			Action:      listProfiles,
			Aliases:     []string{"ls"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
				},
				&cli.StringFlag{
					Name:    "provider",
					Aliases: []string{"prov"},
				},
				&cli.StringFlag{
					Name:    "repo",
					Aliases: []string{"r"},
				},
			},
		},
		&cli.Command{
			Name:        "get-profile",
			Description: "Get profile",
			Action:      getProfile,
			Aliases:     []string{"get"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "source",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name:    "profile",
					Aliases: []string{"prof"},
				},
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"dry"},
				},
			},
		},
		&cli.Command{
			Name:        "diagnostics",
			Action:      runDiagnostics,
			Description: "Run diagnostics",
			Aliases:     []string{"diag"},
		},
		&cli.Command{
			Name:        "install-profile",
			Action:      installProfile,
			Description: "Install profile",
			Aliases:     []string{"install\", \"i"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "source",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name:    "profile",
					Aliases: []string{"prof"},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
