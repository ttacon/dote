package installers // github.com/ttacon/dote/dote/installers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/ttacon/dote/dote/types"
)

func Install(profile *types.Profile) error {
	host, err := sysinfo.Host()
	if err != nil {
		return err
	}

	hinfo := host.Info()
	platform := hinfo.OS.Platform

	var installFn func(profile *types.Profile) error

	switch platform {
	case "darwin":
		installFn = installWithBrew
	case "ubuntu":
		fallthrough
	case "debian":
		installFn = installWithApt
	default:
		return errors.New("unsupported platform")
	}

	return installFn(profile)
}

func installWithBrew(profile *types.Profile) error {
	// NOTE(ttacon): yes, we assume that `brew` is installed.
	for _, tool := range profile.Tools {
		fmt.Println("installing: ", tool.Name)
		cmd := exec.Command("brew", "install", tool.Name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func installWithApt(profile *types.Profile) error {
	for _, tool := range profile.Tools {
		fmt.Println("installing: ", tool.Name)
		cmd := exec.Command("apt", "install", "-y", tool.Name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
