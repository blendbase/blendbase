package gormext

import (
	"database/sql/driver"
	"errors"
	"os"

	"github.com/fernet/fernet-go"
)

type EncryptedValue struct {
	Raw string
}

func (r *EncryptedValue) Scan(value interface{}) error {
	bytes := []byte(value.(string))

	decoded_key, err := fernet.DecodeKeys(os.Getenv("SECRET_ENCRYPTION_KEY"))
	if err != nil {
		return errors.New("bad encryption key")
	}

	decoded_value := fernet.VerifyAndDecrypt(bytes, -1, decoded_key)
	if err != nil {
		return errors.New("failed to decrypt encrypted value")
	}

	r.Raw = string(decoded_value)
	return nil
}

func (r EncryptedValue) Value() (driver.Value, error) {
	encoded_key, err := fernet.DecodeKeys(os.Getenv("SECRET_ENCRYPTION_KEY"))
	if err != nil {
		return nil, errors.New("bad encryption key")
	}

	encrypted, err := fernet.EncryptAndSign([]byte(r.Raw), encoded_key[0])
	if err != nil {
		return nil, errors.New("failed to encrypt value")
	}

	return encrypted, nil
}

func (EncryptedValue) GormDataType() string {
	return "text"
}
