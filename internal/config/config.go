// Package config loads Gateway process configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Hafuunano/Andromeda-Gateway/internal/driver"
)

// Config is process-level Gateway configuration (driver + ingress credentials).
type Config struct {
	Driver driver.Name

	NickNames     []string
	CommandPrefix string
	SuperUsers    []string

	// OneBot
	WSURL   string
	WSToken string

	// QQ Open webhook
	AppID        string
	AppSecret    string
	WebhookHost  string
	WebhookPort  string
	WebhookPath  string
}

// Load reads configuration from the environment. Call after godotenv.Load if desired.
func Load() (Config, error) {
	cfg := Config{
		Driver:        driver.Name(strings.TrimSpace(os.Getenv("DRIVER"))),
		CommandPrefix: envOr("COMMAND_PREFIX", "/"),
		WSURL:         os.Getenv("WS_URL"),
		WSToken:       os.Getenv("WS_TOKEN"),
		AppID:         os.Getenv("APP_ID"),
		AppSecret:     os.Getenv("APP_SECRET"),
		WebhookHost:   envOr("WEBHOOK_HOST", "0.0.0.0"),
		WebhookPort:   envOr("WEBHOOK_PORT", "9000"),
		WebhookPath:   envOr("WEBHOOK_PATH", "/qqbot"),
	}

	if cfg.Driver == "" {
		cfg.Driver = driver.NameOneBot
	}

	for _, p := range strings.Split(os.Getenv("NICK_NAMES"), ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			cfg.NickNames = append(cfg.NickNames, p)
		}
	}
	for _, p := range strings.Split(os.Getenv("SUPER_USERS"), ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			cfg.SuperUsers = append(cfg.SuperUsers, p)
		}
	}

	switch cfg.Driver {
	case driver.NameOneBot, driver.NameQQOpen:
		// ok
	default:
		return Config{}, fmt.Errorf("unsupported DRIVER %q (want onebot|qqopen)", cfg.Driver)
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
