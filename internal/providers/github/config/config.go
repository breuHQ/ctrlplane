package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v62/github"

	"go.breu.io/quantm/internal/db"
	pkgerrors "go.breu.io/quantm/internal/providers/github/errors"
)

type (
	// config holds configuration settings for the GitHub integration.
	config struct {
		AppID              int64  `env:"GITHUB_APP_ID"`                                    // GitHub App ID.
		ClientID           string `env:"GITHUB_CLIENT_ID"`                                 // GitHub Client ID.
		WebhookSecret      string `env:"GITHUB_WEBHOOK_SECRET"`                            // Secret for verifying webhook requests.
		PrivateKey         string `env:"GITHUB_PRIVATE_KEY"`                               // Private key for the GitHub API.
		PrivateKeyIsBase64 bool   `env:"GITHUB_PRIVATE_KEY_IS_BASE64" env-default:"false"` // If true, the private key is base64 encoded.
	}
)

// SignPayload generates a signature for a given payload.
//
// Calculates the HMAC-SHA256 hash of the payload using the webhook secret. Returns the base64 encoded signature in the
// format "sha256=<hash>".
func (cfg *config) SignPayload(payload []byte) string {
	key := hmac.New(sha256.New, []byte(cfg.WebhookSecret))
	key.Write(payload)
	result := "sha256=" + hex.EncodeToString(key.Sum(nil))

	return result
}

// VerifyWebhookSignature verifies the signature of a webhook payload.
//
// Verifies that the provided signature matches the signature generated by signing the payload with the webhook secret.
// Returns an error if the signatures don't match.
func (cfg *config) VerifyWebhookSignature(payload []byte, signature string) error {
	result := cfg.SignPayload(payload)

	if result != signature {
		return pkgerrors.ErrVerifySignature
	}

	return nil
}

// GetClientForInstallationID retrieves a GitHub client for a specific installation ID.
//
// Creates a new client for the specified installation ID using the GitHub Installation API and uses the private key
// from the configuration to authenticate.
func (cfg *config) GetClientForInstallationID(installationID db.Int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, cfg.AppID, installationID.Int64(), []byte(cfg.PrivateKey))
	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})

	return client, nil
}