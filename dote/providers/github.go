package providers // github.com/ttacon/dote/dote/providers

import (
	"context"
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"github.com/ttacon/dote/dote/types"
)

type Provider interface {
	ListProfiles(source string) ([]*types.Profile, error)
	GetProfile(source, profile string) (*types.Profile, error)
}

type GitHubProvider struct {
	client *github.Client
}

func NewGitHubProvider(client *http.Client) Provider {
	return &GitHubProvider{
		client: github.NewClient(client),
	}
}

func (gh *GitHubProvider) ListProfiles(source string) ([]*types.Profile, error) {
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
		profiles []*types.Profile

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
						&types.Profile{
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

func parseGitHubSource(source string) (string, string, string, error) {
	var provider, username, repo string

	sPieces := strings.Split(source, "/")
	if len(sPieces) != 3 {
		return "", "", "", errors.New("invalid source")
	} else if provider, username, repo = sPieces[0], sPieces[1], sPieces[2]; provider != "github.com" {
		return "", "", "", errors.New("invalid provider, must be github.com")
	}

	if provider != "github.com" {
		return "", "", "", errors.New("the GitHub provider only supports github.com sources")
	}
	return provider, username, repo, nil

}

func (gh *GitHubProvider) GetProfile(source, profile string) (*types.Profile, error) {
	_, username, repo, err := parseGitHubSource(source)
	if err != nil {
		return nil, err
	}

	_, dirContents, _, err := gh.client.Repositories.GetContents(
		context.Background(),
		username,
		repo,
		profile,
		&github.RepositoryContentGetOptions{
			Ref: "master", // Make this configurable from the source string later.
		},
	)
	if err != nil {
		return nil, err
	}

	var top = &types.Profile{Name: profile}
	for _, content := range dirContents {
		if *content.Type != "file" {
			continue
		}

		fileContent, _, _, err := gh.client.Repositories.GetContents(
			context.Background(),
			username,
			repo,
			path.Join(profile, content.GetName()),
			&github.RepositoryContentGetOptions{
				Ref: "master", // Make this configurable from the source string later.
			},
		)
		if err != nil {
			return nil, err
		}

		cont, err := fileContent.GetContent()
		if err != nil {
			return nil, err
		}

		var pro types.Profile
		if _, err := toml.Decode(cont, &pro); err != nil {
			return nil, err
		}

		top.Tools = append(top.Tools, pro.Tools...)

		if len(pro.ExtensionProfiles.Profiles) > 0 {
			for _, extensionProfile := range pro.ExtensionProfiles.Profiles {
				top.Extends = append(top.Extends, &types.Profile{
					Name: extensionProfile,
				})
			}
		}
	}

	if len(top.Extends) > 0 {
		for _, extension := range top.Extends {
			extensionPieces := strings.Split(extension.Name, ":")
			if len(extensionPieces) != 2 {
				return nil, errors.New("invalid extension: " + extension.Name)
			}
			source, profileName := extensionPieces[0], extensionPieces[1]

			profile, err := gh.GetProfile(source, profileName)
			if err != nil {
				return nil, err
			}
			extension.Tools = profile.Tools
			extension.Extends = profile.Extends
		}
	}

	return top, nil
}
