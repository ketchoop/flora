package main

import (
	"fmt"
	"log"
	"os"
	"path"

	version "github.com/hashicorp/go-version"
	"github.com/ketchoop/flora"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var (
	VersionNumber     = "Unknown" //nolint:gochecknoglobals
	VersionCommitHash = "Unknown" //nolint:gochecknoglobals
	VersionBuildDate  = "Unknown" //nolint:gochecknoglobals
)

func main() { //nolint:gocyclo
	app := cli.NewApp()
	app.Name = "flora"
	app.Usage = "Simple app to upgrade your terraform"
	app.Version = fmt.Sprintf("%s (%s at %s)", VersionNumber, VersionCommitHash, VersionBuildDate)
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:  "upgrade",
			Usage: "Upgrade terraform",
			Action: func(c *cli.Context) error {
				tfVersion, err := flora.GetLatestVersion()

				if err != nil {
					log.Fatal(err)
				}

				homeDir, _ := homedir.Dir()
				floraPath := path.Join(homeDir, ".flora")

				upgrader := flora.TerraformUpgrader{Version: tfVersion, FloraPath: floraPath}

				if err := upgrader.Run(); err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
		{
			Name:         "download",
			Usage:        "Download specific Terraform version",
			ArgsUsage:    "TERRAFORM_VERSION",
			BashComplete: flora.VersionsCompletion,
			Action: func(c *cli.Context) error {
				versionConstraint := c.Args().First()

				if versionConstraint == "" {
					versionConstraint = flora.GetVersionConstraint()

					if versionConstraint == "" {
						if err := cli.ShowSubcommandHelp(c); err != nil {
							log.Fatal(err)
						}
						return nil
					}
				}

				tfVersion := flora.GetLatestVersionMatchingConstraint(versionConstraint)
				if tfVersion == "" {
					log.Printf("Can't find version matching constraint %s\n", versionConstraint)
					return nil
				}

				homeDir, _ := homedir.Dir()
				floraPath := path.Join(homeDir, ".flora")

				upgrader := flora.TerraformUpgrader{Version: tfVersion, FloraPath: floraPath}

				log.Print("Downloading Terraform " + tfVersion)

				if err := upgrader.DownloadTerraform(); err != nil {
					log.Fatal(err)
				}

				log.Print("Terraform " + tfVersion + " was successfully downloaded")

				return nil
			},
		},
		{
			Name:         "use",
			Usage:        "Download(when it's needed) and use specific terraform version",
			ArgsUsage:    "TERRAFORM_VERSION",
			BashComplete: flora.VersionsCompletion,
			Action: func(c *cli.Context) error {
				versionConstraint := c.Args().First()

				if versionConstraint == "" {
					versionConstraint = flora.GetVersionConstraint()

					if versionConstraint == "" {
						if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
							log.Fatal(err)
						}
						return nil
					}
				}

				tfVersion := flora.GetLatestVersionMatchingConstraint(versionConstraint)
				if tfVersion == "" {
					log.Printf("Can't find version matching constarint %s\n", versionConstraint)
					return nil
				}

				homeDir, _ := homedir.Dir()
				floraPath := path.Join(homeDir, ".flora")

				upgrader := flora.TerraformUpgrader{Version: tfVersion, FloraPath: floraPath}

				if err := upgrader.Run(); err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
		{
			Name:  "versions",
			Usage: "List all available terraform versions",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "num, n",
					Value: 10,
					Usage: "Number of versions print on screen",
				},
				cli.BoolFlag{
					Name:  "local, l",
					Usage: "Show only installed versions of Terraform",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:  "current",
					Usage: "Show currently used version of terraform",
					Action: func(c *cli.Context) error {

						homeDir, _ := homedir.Dir()
						floraPath := path.Join(homeDir, ".flora")

						currentVer, err := flora.GetCurrentVersion(floraPath)

						if err != nil {
							switch err.(type) {
							case *os.PathError:
								fmt.Println("There is no Terraform used(linked) version")
								return nil
							default:
								log.Fatal(err)
							}
						}

						fmt.Printf("Currently used terraform version is %s\n", currentVer)

						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				var versions []*version.Version
				var err error

				homeDir, _ := homedir.Dir()
				floraPath := path.Join(homeDir, ".flora")

				if c.Bool("local") {
					versions, err = flora.ListLocalVersions(floraPath)
				} else {
					versions, err = flora.ListRemoteVersions()
				}

				if err != nil {
					log.Fatal(err)
				}

				if len(versions) == 0 && c.Bool("local") {
					fmt.Printf("There is no packages installed locally\n")

					return nil
				}

				if len(versions) >= c.Int("num") {
					versions = versions[len(versions)-c.Int("num"):]
				}

				curVer, err := flora.GetCurrentVersion(floraPath)

				for _, version := range versions {
					if err == nil && version.Equal(curVer) {
						fmt.Printf("> %s\n", version)
					} else {
						fmt.Printf("  %s\n", version)
					}
				}

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
