package cmd

import (
	"github.com/fernet/fernet-go"
	"github.com/urfave/cli/v2"
)

var GenerateEncryptionKeyCmd = &cli.Command{
	Name:        "gen-enc-key",
	Description: "Use this command to generate a new encryption key",
	Action: func(c *cli.Context) error {
		key := fernet.Key{}
		key.Generate()
		println(key.Encode())
		return nil
	},
}
