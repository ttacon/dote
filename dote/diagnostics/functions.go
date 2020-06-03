package diagnostics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	sysinfo "github.com/elastic/go-sysinfo"
	cli "github.com/urfave/cli/v2"
)

type DiagnosticFunc func(c *cli.Context) error

var (
	DiagnosticFunctions = []DiagnosticFunc{
		GetMachineInformation,
		ListInstalledProfiles,
	}
)

func GetMachineInformation(c *cli.Context) error {
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

func ListInstalledProfiles(c *cli.Context) error {
	fmt.Println("\n\nInstalled profile:")

	dotePath, err := getDotePath()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filepath.Join(
		dotePath,
		"profileIndex.json",
	))
	if err != nil {

		// See if the error is due to the path not existing:
		//
		// - will be os.PathError
		// - Error will be syscall.Errno(0x02)
		pathErr, ok := err.(*os.PathError)
		if !ok {
			return err
		}

		// Check to see if the issue is that the `.config/dote`
		// directory doesn't exist.
		if syscallErno, ok := pathErr.Err.(syscall.Errno); !ok || syscallErno != 2 {
			return err
		}

		// Make the directory
		if err := os.MkdirAll(dotePath, os.ModeDir|0755); err != nil {
			return err
		}

		emptyData := []byte("{}")
		if err := ioutil.WriteFile(filepath.Join(
			dotePath,
			"profileIndex.json",
		), emptyData, 0644); err != nil {
			return err
		}

		// Set it to our newly, empty, index data.
		data = emptyData
	}

	var installedProfiles = map[string]interface{}{}
	if err := json.Unmarshal(data, &installedProfiles); err != nil {
		return err
	}

	if len(installedProfiles) == 0 {
		fmt.Println("no profiles currently installed")
		return nil
	}

	for profileName, _ := range installedProfiles {
		fmt.Println(profileName)
	}

	return nil
}

func getDotePath() (string, error) {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".config", "dote"), nil
}
