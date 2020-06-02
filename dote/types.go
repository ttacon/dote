package main

import (
	"github.com/xlab/treeprint"
)

type Tool struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type Extends struct {
	Profiles []string `toml:"profiles"`
}

type Profile struct {
	Name              string
	ExtensionProfiles Extends `toml:"extends"`
	Extends           []*Profile
	Tools             []Tool `toml:"tool"`
}

func (p *Profile) String() string {
	tree := treeprint.New()
	if len(p.Extends) > 0 {
		extensions := tree.AddBranch("extends")
		for _, extension := range p.Extends {
			extensionBranch := extensions.AddBranch(extension.Name)
			for _, tool := range extension.Tools {
				if len(tool.Version) > 0 {
					extensionBranch.AddMetaNode(tool.Version, tool.Name)
				} else {
					extensionBranch.AddNode(tool.Name)
				}
			}

		}
	}

	for _, tool := range p.Tools {
		if len(tool.Version) > 0 {
			tree.AddMetaNode(tool.Version, tool.Name)
		} else {
			tree.AddNode(tool.Name)
		}
	}

	return tree.String()
}
