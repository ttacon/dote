package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ttacon/dote/dote/diagnostics"
	"github.com/ttacon/dote/dote/installers"
	"github.com/ttacon/dote/dote/providers"
	"github.com/ttacon/dote/dote/storage"
	"github.com/ttacon/dote/dote/types"
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
	client := providers.NewGitHubProvider(nil)
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

func getProfile(c *cli.Context) error {
	source := c.String("source")
	if len(source) == 0 {
		fmt.Println("must provider source")
		os.Exit(1)
	}

	profile := c.String("profile")
	if len(profile) == 0 {
		fmt.Println("must provide profile")
		os.Exit(1)
	}

	// Retrieve profile.
	client := providers.NewGitHubProvider(nil)
	prof, err := client.GetProfile(
		source,
		profile,
	)
	if err != nil {
		fmt.Println("error retrieving profile: ", err)
		os.Exit(1)
	}

	if c.Bool("dry-run") {
		fmt.Println(prof)
		return nil
	}

	return saveProfile(source, profile, prof)
}

func saveProfile(source, profileName string, profile *types.Profile) error {
	strg := storage.NewFSStorage()
	return strg.SaveProfile(source, profileName, profile)
}

func runDiagnostics(c *cli.Context) error {
	strg := storage.NewFSStorage()
	for _, fn := range diagnostics.DiagnosticFunctions {
		if err := fn(c, strg); err != nil {
			fmt.Println("err: ", err)
			return err
		}
	}

	return nil
}

func installProfile(c *cli.Context) error {
	source := c.String("source")
	if len(source) == 0 {
		fmt.Println("must provider source")
		return errors.New("must provider source")
	}

	profileName := c.String("profile")
	if len(profileName) == 0 {
		fmt.Println("must provide profile")
		return errors.New("must provide profile")
	}

	strg := storage.NewFSStorage()
	profile, err := strg.GetProfile(source, profileName)
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	return installers.Install(profile)
}
