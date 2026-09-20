package qqopen

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Hafuunano/Andromeda-Gateway/internal/config"
	"github.com/Hafuunano/Andromeda-Gateway/internal/driver"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol/qqopen"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/interaction/webhook"
	"github.com/tencent-connect/botgo/token"
)

// Driver runs QQ Open Platform Webhook ingress via botgo.
type Driver struct {
	Cfg config.Config
}

// Name implements driver.Driver.
func (d *Driver) Name() driver.Name { return driver.NameQQOpen }

// Capabilities implements driver.Driver.
func (d *Driver) Capabilities() []protocol.Capability {
	return protocol.DefaultQQOpenCapabilities()
}

// Start implements driver.Driver. Blocks on HTTP ListenAndServe until ctx cancel or server error.
func (d *Driver) Start(ctx context.Context, mws []protocol.Middleware) error {
	if d.Cfg.AppID == "" || d.Cfg.AppSecret == "" {
		return fmt.Errorf("qqopen: APP_ID and APP_SECRET are required")
	}

	credentials := &token.QQBotCredentials{
		AppID:     d.Cfg.AppID,
		AppSecret: d.Cfg.AppSecret,
	}
	tokenSource := token.NewQQBotTokenSource(credentials)
	if err := token.StartRefreshAccessToken(ctx, tokenSource); err != nil {
		return fmt.Errorf("qqopen: start token refresh: %w", err)
	}

	api := botgo.NewOpenAPI(credentials.AppID, tokenSource).WithTimeout(5 * time.Second)

	qqopen.InstallWithMiddlewares(qqopen.HostConfig{
		API:           api,
		SuperUsers:    d.Cfg.SuperUsers,
		CommandPrefix: d.Cfg.CommandPrefix,
		Middlewares:   mws,
	})

	mux := http.NewServeMux()
	path := d.Cfg.WebhookPath
	if path == "" {
		path = "/qqbot"
	}
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		webhook.HTTPHandler(w, r, credentials)
	})

	addr := fmt.Sprintf("%s:%s", d.Cfg.WebhookHost, d.Cfg.WebhookPort)
	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
