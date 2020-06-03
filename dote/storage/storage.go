package storage // github.com/ttacon/dote/dote/storage
import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ttacon/dote/dote/types"
)

type Storage interface {
	ListInstalledProfiles() ([]string, error)
	GetProfile(source, profileName string) (*types.Profile, error)
	SaveProfile(source, profileName string, profile *types.Profile) error
}

type fsStorage struct{}

func NewFSStorage() Storage {
	return fsStorage{}
}

var (
	// File names
	profileIndexFileName = "profileIndex.json"
)

func (_ fsStorage) ListInstalledProfiles() ([]string, error) {
	dotePath, err := getDotePath()
	if err != nil {
		return nil, err
	}

	if err := ensureIndexFileExists(dotePath); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(
		dotePath,
		profileIndexFileName,
	))
	if err != nil {
		return nil, err
	}

	var installedProfiles = map[string]interface{}{}
	if err := json.Unmarshal(data, &installedProfiles); err != nil {
		return nil, err
	}

	var names = make([]string, len(installedProfiles))
	var i = 0
	for name, _ := range installedProfiles {
		names[i] = name
		i++
	}

	return names, nil
}

func getDotePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "dote"), nil
}

func (_ fsStorage) SaveProfile(source, profileName string, profile *types.Profile) error {
	dotePath, err := getDotePath()
	if err != nil {
		return err
	} else if err := ensureIndexFileExists(dotePath); err != nil {
		return err
	}

	indexPath := filepath.Join(
		dotePath,
		profileIndexFileName,
	)
	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return err
	}

	var installedProfiles = map[string]interface{}{}
	if err := json.Unmarshal(data, &installedProfiles); err != nil {
		return err
	}

	compoundName := source + ":" + profileName
	if _, ok := installedProfiles[compoundName]; ok {
		// Right now we only support a single version of a profile, so
		// if it's installed, we're done here.
		return nil
	}

	now := time.Now().Unix()
	installedProfiles[compoundName] = interface{}(
		map[string]interface{}{
			"firstInstalledAt": now,
			"updatedAt":        now,
			"versions": []map[string]interface{}{
				map[string]interface{}{
					"version":     "master",
					"installedAt": now,
					"updatedAt":   now,
				},
			},
		},
	)

	savePath := filepath.Join(
		dotePath,
		"profiles",
		compoundName,
		"master",
		"compiledPolicy.json",
	)
	if _, err := os.Stat(savePath); err != nil {
		if err := os.MkdirAll(filepath.Join(
			dotePath,
			"profiles",
			compoundName,
			"master",
		), os.ModeDir|0755); err != nil {
			return err
		}
	}
	profileData, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(savePath, profileData, 0644); err != nil {
		return err
	}

	newData, err := json.Marshal(installedProfiles)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(indexPath, newData, 0644)
}

func ensureIndexFileExists(dotePath string) error {
	_, err := os.Stat(filepath.Join(
		dotePath,
		profileIndexFileName,
	))
	if err == nil {
		return nil
	}

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

	return nil
}

func (_ fsStorage) GetProfile(source, profileName string) (*types.Profile, error) {
	dotePath, err := getDotePath()
	if err != nil {
		return nil, err
	}

	compoundName := source + ":" + profileName

	profilePath := filepath.Join(
		dotePath,
		"profiles",
		compoundName,
		"master",
		"compiledPolicy.json",
	)

	profileRawData, err := ioutil.ReadFile(profilePath)
	if err != nil {
		return nil, err
	}

	var profile types.Profile
	if err := json.Unmarshal(profileRawData, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
