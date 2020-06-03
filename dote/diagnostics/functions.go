package diagnostics

import (
	"fmt"

	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/ttacon/dote/dote/storage"
	cli "github.com/urfave/cli/v2"
)

type DiagnosticFunc func(c *cli.Context, _ storage.Storage) error

var (
	DiagnosticFunctions = []DiagnosticFunc{
		GetMachineInformation,
		ListInstalledProfiles,
	}
)

func GetMachineInformation(c *cli.Context, _ storage.Storage) error {
	fmt.Println("\n\nMachine information:")
	host, err := sysinfo.Host()
	if err != nil {
		return err
	}

	hinfo := host.Info()
	var containerized bool
	if hinfo.Containerized != nil {
		containerized = *hinfo.Containerized
	}
	fmt.Printf(`OS:            %s/%s (%s)
Kernel:        %s
Containerized: %v 
`,
		hinfo.OS.Platform,
		hinfo.OS.Version,
		hinfo.Architecture,
		hinfo.KernelVersion,
		containerized,
	)
	return nil
}

func ListInstalledProfiles(c *cli.Context, strg storage.Storage) error {
	installedProfiles, err := strg.ListInstalledProfiles()
	if err != nil {
		return err
	}

	fmt.Println("\n\nInstalled profiles:")
	if len(installedProfiles) == 0 {
		fmt.Println("no profiles currently installed")
		return nil
	}

	for _, profileName := range installedProfiles {
		fmt.Printf(" - %s\n", profileName)
	}

	return nil
}
