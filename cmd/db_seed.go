package cmd

import (
	"blendbase/config"
	"blendbase/misc/db_utils"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var DBSeedCmd = &cli.Command{
	Name: "db:seed",
	Action: func(c *cli.Context) error {
		godotenv.Load()

		app, err := config.NewApp()
		if err != nil {
			log.Fatalf("failed to create an app: %s", err)
			return err
		}

		err = db_utils.Seed(app)
		if err != nil {
			log.Fatalf("failed to seed the database: %s", err)
			return err
		}

		log.Info("Database seeded successfully")

		return nil
	},
}
