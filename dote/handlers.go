package main

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"
)

func listProfiles(c *cli.Context) error {
	// Check provider support and that username exists
	provider := c.String("provider")
	if len(provider) == 0 {
		fmt.Println("must provide a provider (i.e. \"github.com\"")
		os.Exit(1)
	} else if provider != "github.com" {
		fmt.Println("only github.com is currently supported")
		os.Exit(1)
	}

	username := c.String("username")
	if len(username) == 0 {
		fmt.Println("must provider username")
		os.Exit(1)
	}

	repo := c.String("repo")
	if len(repo) == 0 {
		fmt.Println("must provide repo")
		os.Exit(1)
	}

	// Retrieve repo.
	client := NewGitHubProvider(nil)
	profiles, err := client.ListProfiles(
		fmt.Sprintf(
			"%s/%s/%s",
			provider,
			username,
			repo,
		),
	)
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	fmt.Println("found the following profiles:")
	for _, profile := range profiles {
		fmt.Println(" - ", profile.Name)
	}

	// Parse out profiles.
	//
	// We allow Node style structuring here, so a profile will either be
	// `profile1.toml` or a directory like so:
	//
	//   profile2/
	//     tools.toml
	//     extensions.toml

	return nil
}
