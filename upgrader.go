package flora

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	tfBaseURL       string = "https://releases.hashicorp.com/terraform/%s/terraform_%s.zip"
	tfCheckpointURL string = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
)

type TerraformUpgrader struct {
	Version   string
	floraPath string
}

func InitTerraformUpgrader(version string) *TerraformUpgrader {
	homeDir, _ := homedir.Dir()

	return &TerraformUpgrader{version, homeDir + "/.flora"}
}

func (t TerraformUpgrader) IsDownloadNeeded() bool {
	_, err := os.Stat(t.floraPath + "/terraform_" + t.Version)

	return os.IsNotExist(err)
}

func (t TerraformUpgrader) DownloadTerraform() error {
	tfFileURL := fmt.Sprintf(tfBaseURL, t.Version, t.Version+"_"+runtime.GOOS+"_"+runtime.GOARCH)

	r, err := http.Get(tfFileURL)

	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return errors.New("can't download terraform")
	}

	zipFile, err := os.Create(path.Join(t.floraPath, "terraform_"+t.Version+".zip")) // use pathlib

	if err != nil {
		return err
	}

	defer zipFile.Close()

	defer r.Body.Close()

	_, err = io.Copy(zipFile, r.Body)

	if err != nil {
		return err
	}

	return nil
}

func (t TerraformUpgrader) UnzipAndClean() error {
	_, err := unzip(path.Join(t.floraPath, "terraform_"+t.Version+".zip"), t.floraPath)

	if err != nil {
		return err
	}

	if err = os.Remove(path.Join(t.floraPath, "terraform_"+t.Version+".zip")); err != nil {
		return err
	}

	os.Rename(path.Join(t.floraPath, "terraform"), path.Join(t.floraPath, "terraform_"+t.Version))

	return nil
}

func (t TerraformUpgrader) InstallNewTerraform() error {
	floraBinPath := path.Join(t.floraPath, "bin", "terraform")

	if _, err := os.Lstat(floraBinPath); err == nil {
		os.Remove(floraBinPath)
	}

	log.Print("Adding symlink " + path.Join(t.floraPath, "terraform_"+t.Version) + "->" + floraBinPath)

	if err := os.Symlink(path.Join(t.floraPath, "terraform_"+t.Version), floraBinPath); err != nil {
		return err
	}

	return nil
}

func (t TerraformUpgrader) Run() error {
	fmt.Print(t.IsDownloadNeeded())
	if t.IsDownloadNeeded() {
		log.Print("Downloading Terraform " + t.Version)

		if err := t.DownloadTerraform(); err != nil {
			log.Fatal(err)
		}

		log.Print("Unpacking Terraform " + t.Version)

		if err := t.UnzipAndClean(); err != nil {
			log.Fatal(err)
		}
	}

	if err := t.InstallNewTerraform(); err != nil {
		log.Fatal(err)
	}

	log.Print("Terraform " + t.Version + " was succesfully installed")

	return nil
}
