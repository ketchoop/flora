package flora

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	tfBaseURL       string = "https://releases.hashicorp.com/terraform/%s/terraform_%s.zip"
	tfCheckpointURL string = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
)

type TerraformUpgrader struct {
	Version      string
	tfFileSuffix string // contains version, arch and OS
	floraPath    string
}

func InitTerraformUpgrader(version string) *TerraformUpgrader {
	homeDir, _ := homedir.Dir()

	return &TerraformUpgrader{version, version + "_" + runtime.GOOS + "_" + runtime.GOARCH, homeDir + "/.flora"}
}

func (t TerraformUpgrader) IsDownloadNeeded() bool {
	_, err := os.Stat(t.floraPath + "/terraform_" + t.tfFileSuffix)

	return os.IsNotExist(err)
}

func (t TerraformUpgrader) DownloadTerraform() error {
	tfFileURL := fmt.Sprintf(tfBaseURL, t.Version, t.tfFileSuffix)

	r, err := http.Get(tfFileURL)

	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return errors.New("can't download terraform")
	}

	zipFile, err := os.Create(t.floraPath + "/terraform_" + t.tfFileSuffix + ".zip") // use pathlib

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
	_, err := unzip(t.floraPath+"/terraform_"+t.tfFileSuffix+".zip", t.floraPath) //TODO: use pathlib

	if err != nil {
		return err
	}

	if err = os.Remove(t.floraPath + "/terraform_" + t.tfFileSuffix + ".zip"); err != nil {
		return err
	}

	os.Rename(t.floraPath+"/terraform", t.floraPath+"/terraform_"+t.tfFileSuffix)

	return nil
}

func (t TerraformUpgrader) InstallNewTerraform() error {
	if _, err := os.Lstat(t.floraPath + "/bin/terraform"); err == nil {
		os.Remove(t.floraPath + "/bin/terraform")
	}

	log.Print("Adding symlink " + t.floraPath + "/terraform_" + t.tfFileSuffix + "->" + t.floraPath + "/bin/terraform")

	if err := os.Symlink(t.floraPath+"/terraform_"+t.tfFileSuffix, t.floraPath+"/bin/terraform"); err != nil {
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
