package flora

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	version "github.com/hashicorp/go-version"
)

const releasesIndexURL string = "https://releases.hashicorp.com/terraform/index.json"

func GetLatestVersion() (string, error) {
	type CheckResponse struct {
		CurrentVersion string `json:"current_version"`
	}

	checkResponse := CheckResponse{}

	r, err := http.Get(tfCheckpointURL)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&checkResponse)

	if err != nil {
		return "", err
	}

	return checkResponse.CurrentVersion, nil
}

func GetCurrentVersion() (version string, err error) {
	link, err := os.Readlink(path.Join(floraPath, "bin", "terraform"))

	version = link[strings.LastIndex(link, "_")+1:]

	return
}

func ListLocalVersions() ([]*version.Version, error) {
	var versions []*version.Version
	var rawVersions []string

	tfFilesList, err := filepath.Glob(path.Join(floraPath, "/terraform_*"))

	if err != nil {
		return nil, err
	}

	for _, tfFile := range tfFilesList {
		rawVersions = append(
			rawVersions,
			tfFile[strings.LastIndex(tfFile, "_")+1:],
		)
	}

	versions = make([]*version.Version, len(rawVersions))

	for i, ver := range rawVersions {
		versions[i], _ = version.NewVersion(ver)
	}

	return versions, nil
}

func ListRemoteVersions() ([]*version.Version, error) {
	var versions []*version.Version
	versionsWrapper := struct {
		Versions map[string]interface{}
	}{}

	r, err := http.Get(releasesIndexURL)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&versionsWrapper)

	if err != nil {
		return nil, err
	}

	versions = make([]*version.Version, len(versionsWrapper.Versions))

	i := 0
	for ver := range versionsWrapper.Versions {
		versions[i], _ = version.NewVersion(ver)
		i++
	}

	sort.Sort(version.Collection(versions))

	return versions, nil
}
