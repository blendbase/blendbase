package main

import (
	"blendbase/cmd"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "blendbase"
	app.Description = "Manage you integrations"
	app.HideVersion = true
	app.Commands = []*cli.Command{
		cmd.DBSeedCmd,
		cmd.DBMigrateCmd,
		cmd.ServerCmd,
		cmd.GenerateEncryptionKeyCmd,
		cmd.GenerateAuthTokenCmd,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
