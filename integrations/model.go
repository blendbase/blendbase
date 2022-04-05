package integrations

import (
	"blendbase/misc/gormext"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:UUID;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (base *Base) BeforeCreate(db *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	base.CreatedAt = time.Now()
	base.UpdatedAt = time.Now()
	return nil
}

func (base *Base) BeforeUpdate(db *gorm.DB) error {
	base.UpdatedAt = time.Now()
	return nil
}

type Consumer struct {
	Base
	ConsumerIntegrations []ConsumerIntegration `gorm:"foreignkey:ConsumerID"`
}

type ConsumerOauth2Configuration struct {
	Base
	ConsumerIntegrationID uuid.UUID `gorm:"type:UUID;"`
	ConsumerIntegration   ConsumerIntegration
	ClientID              gormext.EncryptedValue
	ClientSecret          gormext.EncryptedValue
	RedirectURL           string `gorm:"type:VARCHAR(255);"`
	TokenType             string `gorm:"type:VARCHAR(255);"` // e.g. "bearer"
	AccessToken           gormext.EncryptedValue
	RefreshToken          gormext.EncryptedValue
	CustomSettings        datatypes.JSON `gorm:"type:JSONB;"`
}

func (consumerOauth2Configuration *ConsumerOauth2Configuration) GetCustomSettings() (*ConsumerOauth2ConfigurationCustomSettings, error) {
	if consumerOauth2Configuration.CustomSettings == nil {
		return nil, nil
	}

	var customSettings ConsumerOauth2ConfigurationCustomSettings
	if err := json.Unmarshal(consumerOauth2Configuration.CustomSettings, &customSettings); err != nil {
		return nil, err
	}

	return &customSettings, nil
}

type ConsumerOauth2ConfigurationCustomSettings struct {
	SalesforceInstanceSubdomain string `json:"salesforceInstanceSubdomain"`
}

type ConsumerIntegration struct {
	Base
	Type        string    `gorm:"type:VARCHAR(255);"` // e.g. "crm"
	ServiceCode string    `gorm:"type:VARCHAR(255);"` // e.g. "crm_salesforce", "crm_hubspot"
	ConsumerID  uuid.UUID `gorm:"type:UUID;"`
	Enabled     bool      `gorm:"default:false;"`
	Consumer    Consumer
	Secret      gormext.EncryptedValue
}
