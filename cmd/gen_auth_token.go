package cmd

import (
	"os"

	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var consumerId string

var GenerateAuthTokenCmd = &cli.Command{
	Name:        "gen-auth-token",
	Description: "Use this command to generate an authentication token for GraphQL API 'Authorization: Bearer $token'",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "consumer-id",
			Required:    true,
			Destination: &consumerId,
		},
	},
	Action: func(c *cli.Context) error {
		godotenv.Load(".env")

		blendbaseAuthSecret := os.Getenv("BLENDBASE_AUTH_SECRET")
		if blendbaseAuthSecret == "" {
			log.Fatalf("missing BLENDBASE_AUTH_SECRET env var")
		}

		tokenAuth := jwtauth.New("HS256", []byte(blendbaseAuthSecret), nil)

		_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"consumer_id": consumerId})
		if err != nil {
			return err
		}

		println(tokenString)
		return nil
	},
}
