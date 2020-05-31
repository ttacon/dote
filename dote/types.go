package main

type Tool struct {
	Name    string `toml:"name"`
	Version string `toml:"string"`
}

type Extends struct {
	Profiles []string `toml:"profiles"`
}

type Profile struct {
	Name    string
	Extends Extends `toml:"extends"`
	Tools   []Tool  `toml:"tool"`
}
