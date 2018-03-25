package main

import (
	"fmt"
	"log"
	"os"

	version "github.com/hashicorp/go-version"
	"github.com/ketchoop/flora"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "flora"
	app.Usage = "Simple app to upgrade your terraform"
	app.Commands = []cli.Command{
		{
			Name:  "upgrade",
			Usage: "Upgrade terraform",
			Action: func(c *cli.Context) error {
				version, err := flora.GetLatestVersion()

				if err != nil {
					log.Fatal(err)
				}

				upgrader := flora.TerraformUpgrader{version}

				upgrader.Run()

				return nil
			},
		},
		{
			Name:      "use",
			Usage:     "Download(when it's needed) and use specific terraform version",
			ArgsUsage: "TERRAFORM_VERSION",
			Action: func(c *cli.Context) error {
				version := c.Args().First()

				if version == "" {
					cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				}

				upgrader := flora.TerraformUpgrader{version}

				upgrader.Run()

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
			Action: func(c *cli.Context) error {
				var versions []*version.Version
				var err error

				if c.Bool("local") {
					versions, err = flora.ListLocalVersions()
				} else {
					versions, err = flora.ListRemoteVersions()
				}

				if err != nil {
					log.Fatal(err)
				}

				if len(versions) >= c.Int("num") {
					versions = versions[len(versions)-c.Int("num"):]
				}

				for _, version := range versions {
					fmt.Printf("%s\n", version)
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}
