package auth

import (
	"context"
	"duckdb-test/app/internal/config"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
)

var (
	cliName        = "mytool"
	authRecordFile = "auth_record.json"
)

func getAuthPath() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, cliName, authRecordFile)
}

func SaveRecord(record azidentity.AuthenticationRecord) error {
	path := getAuthPath()
	os.MkdirAll(filepath.Dir(path), 0700)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(record)
}

func LoadRecord() (azidentity.AuthenticationRecord, error) {
	path := getAuthPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return azidentity.AuthenticationRecord{}, err
	}
	var record azidentity.AuthenticationRecord
	err = json.Unmarshal(data, &record)
	return record, err
}

func AuthenticateAndSave(cfg *config.Config) error {

	// Try to load existing record
	if _, err := GetAuthToken(cfg); err == nil {
		fmt.Println("✅ Already logged in.")
		return nil
	}

	// No record or expired/missing token — do interactive login
	cred, err := newCredential(nil, cfg)
	if err != nil {
		return err
	}

	p := policy.TokenRequestOptions{
		Scopes: cfg.Scopes,
	}

	record, err := cred.Authenticate(context.Background(), &p)

	if err != nil {
		return err
	}

	if err := SaveRecord(record); err != nil {
		return err
	}

	fmt.Println("✅ Logged in via browser and saved session.")
	return nil
}

func GetAuthToken(cfg *config.Config) (*string, error) {

	if record, err := LoadRecord(); err == nil {

		cred, err := newCredential(&record, cfg)
		if err != nil {
			return nil, err
		}

		// Try to get a token silently from cache
		token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
			TenantID: cfg.TenantID,
			Scopes:   cfg.Scopes,
		})

		if err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}

		return &token.Token, nil
	}

	return nil, fmt.Errorf("authentication failed: no saved session found")
}

func Logout() error {
	// Remove the saved authentication record
	path := getAuthPath()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove auth record: %w", err)
	}

	fmt.Println("👋 Logged out successfully.")
	return nil
}

func newCredential(authRecord *azidentity.AuthenticationRecord, cfg *config.Config) (*azidentity.InteractiveBrowserCredential, error) {
	c, err := cache.New(nil)
	if err != nil {
		return nil, err
	}

	opts := azidentity.InteractiveBrowserCredentialOptions{
		ClientID: cfg.ClientID,
		TenantID: cfg.TenantID,
		Cache:    c,
	}

	// Only set the authentication record if one was provided.
	if authRecord != nil {
		opts.AuthenticationRecord = *authRecord
	}

	cred, err := azidentity.NewInteractiveBrowserCredential(&opts)
	if err != nil {
		return nil, err
	}
	return cred, nil
}
