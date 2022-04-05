package cmd

import (
	"blendbase/connectors/salesforce"
	"blendbase/graph"
	"blendbase/graph/auth"
	"blendbase/graph/generated"
	"blendbase/integrations"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"blendbase/config"

	// "io/ioutil"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

const defaultPort = "8080"

var (
	app *config.App
)

func ConsumerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumerID := chi.URLParam(r, "consumerID")

		consumer := integrations.Consumer{}

		if err := app.DB.Where("id = ?", consumerID).First(&consumer).Error; err != nil {
			errorMessage := fmt.Sprintf("Error finding consumer %s: %s", consumerID, err)
			log.Error(errorMessage)
			http.Error(w, errorMessage, 404)

			return
		}

		ctx := context.WithValue(r.Context(), "consumerID", consumerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var ServerCmd = &cli.Command{
	Name: "server",
	Action: func(c *cli.Context) error {
		godotenv.Load(".env")
		app, _ = config.NewApp()

		graphAuth := auth.NewAuth()

		port := os.Getenv("PORT")
		if port == "" {
			port = defaultPort
		}

		omniAPIServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
			App:       app,
			GraphAuth: graphAuth,
		}}))

		r := chi.NewRouter()

		// CORS
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-api-token"},
			// ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
			Debug:            true,
		}))

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("OK"))
		})

		r.Route("/connect", func(r chi.Router) {
			r.Route("/{consumerID:[0-9a-f-]+}/integrations", func(r chi.Router) {
				r.Use(ConsumerCtx)

				r.Route("/crm_salesforce/oauth2", func(r chi.Router) {
					r.Get("/login", salesforce.AuthHandleLogin(app))
					r.Get("/callback", salesforce.AuthHandleCallback(app))
				})
			})
		})

		// Omni API
		r.Route("/omni", func(r chi.Router) {
			r.Use(graphAuth.Verifier())
			r.Use(graphAuth.Authenticator)
			r.Handle("/query", omniAPIServer)
		})

		// GrahQL playground
		// TODO: Dev env only
		r.Handle("/", playground.Handler("GraphQL playground", "/omni/query"))

		app.Logger.Infof("Connect to http://localhost:%s/ for GraphQL playground", port)
		app.Logger.Fatal(http.ListenAndServe(":"+port, r))

		return nil
	},
}
