package salesforce

import (
	"blendbase/config"
	"blendbase/connectors"
	"blendbase/integrations"
	"blendbase/misc/gormext"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2"
)

func LoadConsumerFromRequestContext(app *config.App, r *http.Request) (*integrations.Consumer, error) {
	ctx := r.Context()
	consumerID, ok := ctx.Value("consumerID").(string)
	consumer := integrations.Consumer{}
	if !ok {
		return nil, errors.New("Missing consumer ID")
	}

	if err := app.DB.Where("id = ?", consumerID).First(&consumer).Error; err != nil {
		return nil, err
	}

	return &consumer, nil
}

func AuthHandleLogin(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		consumer, err := LoadConsumerFromRequestContext(app, r)
		if err != nil {
			errorMessage := fmt.Sprintf("Error finding consumer: %s", err)
			http.Error(w, errorMessage, http.StatusUnprocessableEntity)
			return
		}

		client, err := LoadClientFromDB(app, consumer)
		if err != nil {
			errorMessage := fmt.Sprintf("Error loading consumer: %s", err)
			log.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusUnprocessableEntity)
			return
		}

		url := client.GetAuthCodeUrl()
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func AuthHandleCallback(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const ServiceType = connectors.CONNECTOR_CRM_SALESFORCE
		clientIntegrationsPageURL := os.Getenv("CLIENT_APP_INTEGRATIONS_PAGE_URL")

		consumer, err := LoadConsumerFromRequestContext(app, r)
		if err != nil {
			errorMessage := fmt.Sprintf("Error finding consumer: %s", err)
			http.Error(w, errorMessage, http.StatusUnprocessableEntity)
			return
		}

		client, err := LoadClientFromDB(app, consumer)
		if err != nil {
			errorMessage := fmt.Sprintf("Error loading consumer: %s", err)
			log.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusUnprocessableEntity)
			return
		}

		token, err := client.GetToken(r.FormValue("state"), r.FormValue("code"))
		if err != nil {
			app.Logger.Errorf("Error getting OAuth token from Salesforce: %s", err)
			url := fmt.Sprintf("%s?blendbaseErrorMessage=%s", clientIntegrationsPageURL, "Error getting OAuth token from Salesforce.")
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}

		// Look for existing salesforce integration
		var consumerIntegration integrations.ConsumerIntegration

		if err := app.DB.Where("consumer_id = ?", consumer.ID).Where("service_code = ?", ServiceType).First(&consumerIntegration).Error; err != nil {
			errorString := fmt.Sprintf("Error finding %s integration for %s: %s", ServiceType, consumer.ID, err)
			log.Error(errorString)

			url := fmt.Sprintf("%s?blendbaseErrorMessage=%s", clientIntegrationsPageURL, "Error finding consumer ID.")
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}

		updatedAuthConfig := integrations.ConsumerOauth2Configuration{
			ConsumerIntegrationID: consumerIntegration.ID,
			TokenType:             token.TokenType,
			AccessToken: gormext.EncryptedValue{
				Raw: token.AccessToken,
			},
			RefreshToken: gormext.EncryptedValue{
				Raw: token.RefreshToken,
			},
		}

		// Look for existing salesforce OAuth2 configuration
		var oauthConfig integrations.ConsumerOauth2Configuration
		if err := app.DB.Where("consumer_integration_id = ?", consumerIntegration.ID).First(&oauthConfig).Error; err != nil {
			app.Logger.Infof("%s OAuth2 configuration does not exist, creating…", ServiceType)

			if err := app.DB.Create(&oauthConfig).Error; err != nil {
				errorMessage := fmt.Sprintf("Error creating %s OAuth2 configuration: %s", ServiceType, err)
				log.Error(errorMessage)

				url := fmt.Sprintf("%s?blendbaseErrorMessage=%s", clientIntegrationsPageURL, "Error finding the integration.")
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			}
		}

		log.Infof("%s OAuth2 configuration exists, updating…", ServiceType)
		if err := app.DB.Model(&oauthConfig).Where("consumer_integration_id = ?", consumerIntegration.ID).Updates(updatedAuthConfig).Error; err != nil {
			errorMessage := fmt.Sprintf("Error updating %s OAuth2 configuration: %s", ServiceType, err)
			app.Logger.Error(errorMessage)

			url := fmt.Sprintf("%s?blendbaseErrorMessage=%s", clientIntegrationsPageURL, "Error updating OAuth2 token. Please try again.")
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}
		app.Logger.Infof("%s OAuth2 configuration updated", ServiceType)

		url := fmt.Sprintf("%s?blendbaseSuccessMessage=%s", clientIntegrationsPageURL, "Salesforce OAuth2 token was updated")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (client *Client) GetToken(state string, code string) (*oauth2.Token, error) {
	if state != client.OAuthStateString {
		log.Fatal("invalid oauth state")
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := getOAuthConfig(client.consumerOAuthConfig).Exchange(context.Background(), code)
	if err != nil {
		errorMessage := fmt.Sprintf("code exchange failed: %s", err.Error())
		log.Fatal(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}

	return token, nil
}

func (client *Client) refreshToken() error {
	log.Info("SF refreshing token")

	config := getOAuthConfig(client.consumerOAuthConfig)

	// omitting AccessToken to force a refresh
	expiredToken := oauth2.Token{
		RefreshToken: client.consumerOAuthConfig.RefreshToken.Raw,
		TokenType:    client.consumerOAuthConfig.TokenType,
	}

	newToken, err := config.TokenSource(context.TODO(), &expiredToken).Token()

	if err != nil {
		return errors.New("failed to refresh token")
	}

	err = client.app.DB.Model(client.consumerOAuthConfig).Updates(integrations.ConsumerOauth2Configuration{
		AccessToken: gormext.EncryptedValue{
			Raw: newToken.AccessToken,
		},
		RefreshToken: gormext.EncryptedValue{
			Raw: newToken.RefreshToken,
		},
	}).Error

	client.HTTPClient = config.Client(context.TODO(), newToken)

	return err
}

func (client *Client) GetAuthCodeUrl() string {
	return getOAuthConfig(client.consumerOAuthConfig).AuthCodeURL(client.OAuthStateString)
}

func getOAuthConfig(consumerOAuthConfig *integrations.ConsumerOauth2Configuration) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  consumerOAuthConfig.RedirectURL,
		ClientID:     consumerOAuthConfig.ClientID.Raw,
		ClientSecret: consumerOAuthConfig.ClientSecret.Raw,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.salesforce.com/services/oauth2/authorize",
			TokenURL: "https://login.salesforce.com/services/oauth2/token",
		},
	}
}

func getOAuthToken(consumerOAuthConfig *integrations.ConsumerOauth2Configuration) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  consumerOAuthConfig.AccessToken.Raw,
		RefreshToken: consumerOAuthConfig.RefreshToken.Raw,
		TokenType:    consumerOAuthConfig.TokenType,
	}
}
