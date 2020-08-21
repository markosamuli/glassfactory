package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/markosamuli/glassfactory"
	"github.com/markosamuli/glassfactory/api"
	"github.com/pkg/errors"
	"github.com/zalando/go-keyring"
)

// Auth represents the authentication details for Glass Factory
type Auth struct {
	Account string
	Email   string
	APIKey  string
}

type key int

const (
	gfAuthKey key = iota
)

// NewContext adds Auth to the current context
func NewContext(ctx context.Context, client *Auth) context.Context {
	return context.WithValue(ctx, gfAuthKey, client)
}

// FromContext returns Auth from the current context
func FromContext(ctx context.Context) (*Auth, bool) {
	auth, ok := ctx.Value(gfAuthKey).(*Auth)
	return auth, ok
}

// NewAuth creates a new authentication details
func NewAuth() *Auth {
	return &Auth{}
}

// Setup loads authentication details
func (b *Auth) Setup() error {
	if b.Account == "" {
		account := os.Getenv("GF_ACCOUNT")
		if account != "" {
			b.Account = account
		} else {
			return fmt.Errorf("missing Glass Factory account subdomain")
		}
	}
	if b.Email == "" {
		email := os.Getenv("GF_EMAIL")
		if email != "" {
			b.Email = email
		} else {
			return fmt.Errorf("missing Glass Factory user email address")
		}
	}
	if b.APIKey == "" {
		apiKey := os.Getenv("GF_API_KEY")
		if apiKey != "" {
			b.APIKey = apiKey
		}
	}
	if b.APIKey == "" {
		keyringService, err := b.keyringService()
		if err != nil {
			return err
		}
		if apiKey, err := keyring.Get(keyringService, b.Email); err == nil {
			b.APIKey = apiKey
		} else {
			return errors.Wrapf(err, "failed to get Glass Factory login details for user %s", b.Email)
		}
	}
	return nil
}

// DeleteLoginDetailsFromKeyring deletes username and password from the keyring
func (b *Auth) DeleteLoginDetailsFromKeyring() error {
	keyringService, err := b.keyringService()
	if err != nil {
		return err
	}
	return keyring.Delete(keyringService, b.Email)
}

// StoreLoginDetailsInKeyring stores username and password in the keyring
func (b *Auth) StoreLoginDetailsInKeyring() error {
	keyringService, err := b.keyringService()
	if err != nil {
		return err
	}
	return keyring.Set(keyringService, b.Email, b.APIKey)
}

func (b *Auth) keyringService() (string, error) {
	if b.Account == "" {
		return "", fmt.Errorf("missing Glass Factory account subdomain")
	}
	return fmt.Sprintf("%s.glassfactory.io", b.Account), nil

}

// Validate returns an error if any authentication settings are missing
func (b *Auth) Validate() error {
	if b.Account == "" || b.Email == "" || b.APIKey == "" {
		return errors.New("missing Glass Factory account or authentication details")
	}
	return nil
}

// NewService creates new Glass Factory service
func (b *Auth) NewService() (*api.Service, error) {
	if b.Account == "" || b.Email == "" || b.APIKey == "" {
		err := b.Setup()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new Glass Factory service")
		}
	}
	service, err := glassfactory.New(b.Account, b.Email, b.APIKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new Glass Factory service")
	}
	return service, nil
}
