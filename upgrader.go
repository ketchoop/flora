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
)

const (
	tfBaseURL       string = "https://releases.hashicorp.com/terraform/%s/terraform_%s.zip"
	tfCheckpointURL string = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
)

type TerraformUpgrader struct {
	Version   string
	FloraPath string
}

func (t TerraformUpgrader) IsDownloadNeeded() bool {
	_, err := os.Stat(t.FloraPath + "/terraform_" + t.Version)

	return os.IsNotExist(err)
}

func (t TerraformUpgrader) DownloadTerraform() error {
	tfFileURL := fmt.Sprintf(tfBaseURL, t.Version, t.Version+"_"+runtime.GOOS+"_"+runtime.GOARCH)

	r, err := http.Get(tfFileURL) //nolint:gosec

	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return errors.New("can't download terraform")
	}

	zipFile, err := os.Create(path.Join(t.FloraPath, "terraform_"+t.Version+".zip")) // use pathlib

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
	_, err := unzip(path.Join(t.FloraPath, "terraform_"+t.Version+".zip"), t.FloraPath)

	if err != nil {
		return err
	}

	if err := os.Remove(path.Join(t.FloraPath, "terraform_"+t.Version+".zip")); err != nil {
		return err
	}

	if err := os.Rename(path.Join(t.FloraPath, "terraform"), path.Join(t.FloraPath, "terraform_"+t.Version)); err != nil {
		return err
	}

	return nil
}

func (t TerraformUpgrader) InstallNewTerraform() error {
	floraBinPath := path.Join(t.FloraPath, "bin", "terraform")

	if _, err := os.Lstat(floraBinPath); err == nil {
		os.Remove(floraBinPath)
	}

	log.Print("Adding symlink " + path.Join(t.FloraPath, "terraform_"+t.Version) + "->" + floraBinPath)

	if err := os.Symlink(path.Join(t.FloraPath, "terraform_"+t.Version), floraBinPath); err != nil {
		return err
	}

	return nil
}

func (t TerraformUpgrader) Run() error {
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

	log.Print("Terraform " + t.Version + " was successfully installed")

	return nil
}
