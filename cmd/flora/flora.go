package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ketchoop/flora"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "flora"
	app.Usage = "Simple app to upgrade your terraform"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "upgrade",
			Usage: "Upgrade terraform",
			Action: func(c *cli.Context) error {
				version, err := flora.GetLatestVersion()

				if err != nil {
					log.Fatal(err)
				}

				upgrader := flora.InitTerraformUpgrader(version)

				upgrader.Run(time.Now())

				return nil
			},
		},
		{
			Name:      "install",
			Usage:     "Install specific terraform version",
			ArgsUsage: "TERRAFORM_VERSION",
			Action: func(c *cli.Context) error {
				version := c.Args().First()

				if version == "" {
					cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				}

				upgrader := flora.InitTerraformUpgrader(version)

				upgrader.Run(time.Now())

				return nil
			},
		},
		{
			Name:  "versions",
			Usage: "List all available terraform versions",
			Action: func(c *cli.Context) error {
				versions, err := flora.ListAllVersions()

				if err != nil {
					log.Fatal(err)
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
