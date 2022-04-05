package auth

import (
	"blendbase/misc/db_utils"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	jwt "github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
)

type contextKey struct {
	name string
}

type GraphAuth struct {
	tokenAuth  *jwtauth.JWTAuth
	contextKey *contextKey
}

func NewAuth() *GraphAuth {
	blendbaseAuthSecret := os.Getenv("BLENDBASE_AUTH_SECRET")
	if blendbaseAuthSecret == "" {
		log.Fatalf("missing BLENDBASE_AUTH_SECRET env var")
	}

	return &GraphAuth{
		tokenAuth:  jwtauth.New("HS256", []byte(blendbaseAuthSecret), nil),
		contextKey: &contextKey{"consumer_id"},
	}
}

func (graphAuth *GraphAuth) GenerateTestToken() (string, error) {
	_, tokenString, err := graphAuth.tokenAuth.Encode(map[string]interface{}{"consumer_id": db_utils.TEST_CONSUMER_ID})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (graphAuth *GraphAuth) Verifier() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verify(graphAuth.tokenAuth, jwtauth.TokenFromHeader)(next)
	}
}

func (graphAuth *GraphAuth) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil || token == nil || jwt.Validate(token) != nil {
			http.Error(w, formatError("Invalid or missing JWT token"), http.StatusUnauthorized)
			return
		}

		if len(claims) > 0 && claims[graphAuth.contextKey.name] != nil {
			consumerID := claims[graphAuth.contextKey.name].(string)
			consumerIDuuid, err := uuid.Parse(consumerID)
			if err != nil {
				http.Error(w, formatError("Unprocessable consumer ID claim. Please provide a valid UUID consumer ID"), http.StatusUnprocessableEntity)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), graphAuth.contextKey, consumerIDuuid))
		}
		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func (graphAuth *GraphAuth) GetConsumerIDFromContext(ctx context.Context) *uuid.UUID {
	consumerID := ctx.Value(graphAuth.contextKey)

	if consumerID == nil {
		return nil
	}
	consumerIDuuid := consumerID.(uuid.UUID)

	return &consumerIDuuid
}

func formatError(errorMessage string) string {
	if errorMessage == "" {
		return "{}"
	}
	return fmt.Sprintf(`{"message": "%s"}`, errorMessage)
}
