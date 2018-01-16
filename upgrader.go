package flora

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	tfBaseURL       string = "https://releases.hashicorp.com/terraform/%s/terraform_%s.zip"
	tfDownloadPath  string = "/tmp"
	tfCheckpointURL string = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
)

type TerraformUpgrader struct {
	Version      string
	tfFileSuffix string // contains version, arch and OS
}

func InitTerraformUpgrader(version string) *TerraformUpgrader {
	return &TerraformUpgrader{version, version + "_" + runtime.GOOS + "_" + runtime.GOARCH}
}

//func (t TerraformUpgrader) IsUpgradeNeeded() {
//	oldTfVersion, err = os.Exec()
//
//	if err != nil {
//		return err
//	}
//
//	return t.Version == oldTfVersion
//}
func timeTrack(start time.Time, action string) {
	elapsed := time.Since(start)
	fmt.Printf("Terraform %s in %s\n\n", action, elapsed)
}

func (t TerraformUpgrader) DownloadTerraform() error {
	tfFileURL := fmt.Sprintf(tfBaseURL, t.Version, t.tfFileSuffix)

	r, err := http.Get(tfFileURL)
	// fmt.Println(r)

	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.New("can't download terraform")
	}

	zipFile, err := os.Create(tfDownloadPath + "/terraform_" + t.tfFileSuffix + ".zip") // use pathlib

	if err != nil {
		return err
	}

	defer zipFile.Close()


	_, err = io.Copy(zipFile, r.Body)

	if err != nil {
		return err
	}

	return nil
}

func (t TerraformUpgrader) UnzipAndClean() error {
	_, err := unzip(tfDownloadPath+"/terraform_"+t.tfFileSuffix+".zip", tfDownloadPath) // use pathlib

	if err != nil {
		return err
	}

	return nil
}

func (t TerraformUpgrader) InstallNewTerraform() error {
	oldTfPath, err := exec.LookPath("terraform")

	if err != nil {
		oldTfPath = "/usr/bin/terraform"
	}

	err = os.Rename(tfDownloadPath+"/terraform", oldTfPath)

	if err != nil {
		return err
	}

	return nil
}

func (t TerraformUpgrader) Run(startTime time.Time) error {
	fmt.Printf("Step 1/3: Downloading Terraform %s to %s\n", t.Version, tfDownloadPath)
	if err := t.DownloadTerraform(); err != nil {
		log.Fatal(err)
	}
	timeTrack(startTime, "Downloaded")

	fmt.Printf("Step 2/3: Unpacking Terraform %s to %s\n", t.Version, tfDownloadPath)

	startTime = time.Now()
	if err := t.UnzipAndClean(); err != nil {
		log.Fatal(err)
	}
	timeTrack(startTime, "Unpacked")

	fmt.Println("Step 3/3: Installing Terraform to /usr/bin/")

	startTime = time.Now()
	if err := t.InstallNewTerraform(); err != nil {
		log.Fatal(err)
	}
	timeTrack(startTime, "Installed")

	return nil
}
