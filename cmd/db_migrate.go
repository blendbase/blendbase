package cmd

import (
	"blendbase/config"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"blendbase/misc/db_utils"
)

var DBMigrateCmd = &cli.Command{
	Name: "db:migrate",
	Action: func(c *cli.Context) error {
		godotenv.Load()

		app, err := config.NewApp()
		if err != nil {
			app.Logger.Fatal("Failed to create an app")
			return err
		}

		err = db_utils.Migrate(app)

		if err != nil {
			app.Logger.Fatalf("Failed to migrate database: %s", err)
			return err
		}

		app.Logger.Info("Successfully migrated the database")

		return nil
	},
}
