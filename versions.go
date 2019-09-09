package flora

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

func GetCurrentVersion(floraPath string) (ver *version.Version, err error) {
	link, err := os.Readlink(path.Join(floraPath, "bin", "terraform"))

	if err != nil {
		return
	}

	ver, err = version.NewVersion(link[strings.LastIndex(link, "_")+1:])

	return
}

func ListLocalVersions(floraPath string) ([]*version.Version, error) {
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

func getVersionConstraintFromFile(file string) string {
	f, err := os.Open(file)

	if err != nil {
		return ""
	}
	defer f.Close()

	re := regexp.MustCompile(`required_version[ \t]*=[ \t]*"(.*)"`)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		str := scanner.Text()
		result := re.FindStringSubmatch(str)
		if len(result) >= 2 && result[1] != "" {
			return result[1]
		}
	}

	return ""
}

func GetVersionConstraint() string {
	tfFiles, _ := filepath.Glob("*.tf")
	for _, file := range tfFiles {
		versionConstrains := getVersionConstraintFromFile(file)
		if versionConstrains != "" {
			return versionConstrains
		}
	}

	return ""
}

func getVersionMatchingConstraint(constraintString string, versions []*version.Version) *version.Version {
	constraint, _ := version.NewConstraint(constraintString)
	for i := len(versions) - 1; i >= 0; i-- {
		if constraint.Check(versions[i]) {
			return versions[i]
		}
	}
	return nil
}

func GetLatestVersionMatchingConstraint(versionConstraint string) string {
	versions, _ := ListRemoteVersions()
	tfVersion := getVersionMatchingConstraint(versionConstraint, versions)
	if tfVersion == nil {
		return ""
	}

	return tfVersion.String()
}
