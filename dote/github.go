package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
)

type Provider interface {
	ListProfiles(source string) ([]*Profile, error)
	GetProfile(source, profile string) (*Profile, error)
}

type GitHubProvider struct {
	client *github.Client
}

func NewGitHubProvider(client *http.Client) Provider {
	return &GitHubProvider{
		client: github.NewClient(client),
	}
}

func (gh *GitHubProvider) ListProfiles(source string) ([]*Profile, error) {
	var provider, username, repo string

	sPieces := strings.Split(source, "/")
	if len(sPieces) != 3 {
		return nil, errors.New("invalid source")
	} else if provider, username, repo = sPieces[0], sPieces[1], sPieces[2]; provider != "github.com" {
		return nil, errors.New("invalid provider, must be github.com")
	}

	if provider != "github.com" {
		return nil, errors.New("the GitHub provider only supports github.com sources")
	}

	var (
		profiles []*Profile

		unseenPath  string
		unseenPaths = []string{""}
	)

	for len(unseenPaths) > 0 {
		unseenPath, unseenPaths = unseenPaths[0], unseenPaths[1:]

		_, dirContents, _, err := gh.client.Repositories.GetContents(
			context.Background(),
			username,
			repo,
			unseenPath,
			&github.RepositoryContentGetOptions{
				Ref: "master", // Make this configurable from the source string later.
			},
		)
		if err != nil {
			return nil, err
		}

		foundProfile := false

		for _, content := range dirContents {
			// Special case, top-level profiles.
			if *content.Type == "file" && unseenPath != "" {
				if !foundProfile && strings.HasSuffix(
					*content.Name,
					".toml",
				) {
					foundProfile = true
					profiles = append(
						profiles,
						&Profile{
							Name: unseenPath,
						},
					)
				}
				continue
			}

			// Assert it's a dir and look for at least one toml file.
			if *content.Type == "dir" {
				unseenPaths = append(unseenPaths, *content.Path)
			}
		}
	}
	return profiles, nil
}

func (gh *GitHubProvider) GetProfile(source, profile string) (*Profile, error) {
	return nil, errors.New("not implemented")
}
