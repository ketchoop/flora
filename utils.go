package flora

import (
	"fmt"

	"github.com/urfave/cli"
)

func VersionsCompletion(c *cli.Context) {
	if c.NArg() > 0 {
		return
	}

	versions, err := ListRemoteVersions()

	if err != nil {
		return
	}

	for _, version := range versions {
		fmt.Println(version)
	}
}
