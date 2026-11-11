package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyOAuthEnvOverrides(t *testing.T) {
	c := &Config{}

	os.Setenv("GOOGLE_CLIENT_ID", "g-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "g-secret")
	os.Setenv("GITHUB_CLIENT_ID", "gh-id")
	os.Setenv("GITHUB_CLIENT_SECRET", "gh-secret")
	defer func() {
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")
		os.Unsetenv("GITHUB_CLIENT_ID")
		os.Unsetenv("GITHUB_CLIENT_SECRET")
	}()

	applyOAuthEnvOverrides(c)

	if c.OAuth.Google.ClientID != "g-id" {
		t.Errorf("google client id = %q, want %q", c.OAuth.Google.ClientID, "g-id")
	}
	if c.OAuth.Google.ClientSecret != "g-secret" {
		t.Errorf("google client secret = %q, want %q", c.OAuth.Google.ClientSecret, "g-secret")
	}
	if c.OAuth.Github.ClientID != "gh-id" {
		t.Errorf("github client id = %q, want %q", c.OAuth.Github.ClientID, "gh-id")
	}
	if c.OAuth.Github.ClientSecret != "gh-secret" {
		t.Errorf("github client secret = %q, want %q", c.OAuth.Github.ClientSecret, "gh-secret")
	}
}

func TestApplyOAuthEnvOverridesEmptyDoesNotClobber(t *testing.T) {
	c := &Config{}
	c.OAuth.Google.ClientID = "file-id"

	applyOAuthEnvOverrides(c)

	if c.OAuth.Google.ClientID != "file-id" {
		t.Errorf("empty env var clobbered existing value: got %q", c.OAuth.Google.ClientID)
	}
}

func TestLoadConfigReadsOAuthFromEnv(t *testing.T) {
	os.Setenv("GOOGLE_CLIENT_SECRET", "from-env")
	defer os.Unsetenv("GOOGLE_CLIENT_SECRET")

	path := filepath.Join(t.TempDir(), "configs.json")
	content := `{
		"frontend": {"url": "http://localhost:3000"},
		"oauth": {
			"google": {"client_id": "", "client_secret": ""},
			"github": {"client_id": "", "client_secret": ""}
		}
	}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.OAuth.Google.ClientSecret != "from-env" {
		t.Errorf("google secret = %q, want %q", cfg.OAuth.Google.ClientSecret, "from-env")
	}
}
