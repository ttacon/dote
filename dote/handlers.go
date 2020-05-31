package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/kr/pretty"
	"github.com/urfave/cli"
)

func listProfiles(c *cli.Context) {
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
	client := github.NewClient(nil)

	_,
		dirContents,
		_,
		err := client.Repositories.GetContents(
		context.Background(),
		username,
		repo,
		"",
		&github.RepositoryContentGetOptions{
			Ref: "master",
		},
	)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
		} else {
			fmt.Println("err: ", err)
			os.Exit(1)
		}
	}

	var profiles []string

	for _, content := range dirContents {
		pretty.Println(content)
		name := *content.Name
		if *content.Type == "file" {
			if strings.HasSuffix(
				*content.Name,
				".toml",
			) {
				profiles = append(
					profiles,
					strings.TrimSuffix(
						name,
						".toml",
					),
				)
			}
			continue
		}

		// Assert it's a dir and look for at least one toml file.
		if *content.Type != "dir" {
			continue
		}

		_, nestedDir, _, err := client.Repositories.GetContents(
			context.Background(),
			username,
			repo,
			*content.Path,
			&github.RepositoryContentGetOptions{
				Ref: "master",
			},
		)
		if err != nil {
			fmt.Println("err: ", err)
			os.Exit(1)
		}
		foundTomlFile := false
		for _, nested := range nestedDir {
			// For now, only support one level of nesting.
			if *nested.Type != "file" {
				continue
			}

			if strings.HasSuffix(
				*nested.Name,
				".toml",
			) {
				foundTomlFile = true
				break
			}
		}
		if foundTomlFile {
			profiles = append(
				profiles,
				name,
			)
		}
	}

	fmt.Println("found the following profiles:")
	for _, profile := range profiles {
		fmt.Println(" - ", profile)
	}

	// Parse out profiles.
	//
	// We allow Node style structuring here, so a profile will either be
	// `profile1.toml` or a directory like so:
	//
	//   profile2/
	//     tools.toml
	//     extensions.toml

}
