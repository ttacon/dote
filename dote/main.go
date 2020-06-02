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
			Action:      listProfiles,
			Name:        "list-profiles",
			Description: "List profiles",
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
			Action:      getProfile,
			Name:        "get-profile",
			Description: "Get profile",
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
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
