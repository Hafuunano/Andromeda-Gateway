package sandbox

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Hafuunano/Andromeda-Gateway/internal/config"
	"github.com/Hafuunano/Andromeda-Gateway/internal/driver"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	pcsandbox "github.com/Hafuunano/Protocol-ConvertTool/protocol/sandbox"
)

// Driver serves the sandbox IM web UI and WebSocket protocol.
type Driver struct {
	Cfg config.Config
}

// Name implements driver.Driver.
func (d *Driver) Name() driver.Name { return driver.NameSandbox }

// Capabilities implements driver.Driver.
func (d *Driver) Capabilities() []protocol.Capability {
	return protocol.DefaultSandboxCapabilities()
}

// Start implements driver.Driver.
func (d *Driver) Start(ctx context.Context, mws []protocol.Middleware) error {
	mux := http.NewServeMux()
	pcsandbox.Mount(mux, pcsandbox.HostConfig{
		SuperUsers:    d.Cfg.SuperUsers,
		NickNames:     d.Cfg.NickNames,
		CommandPrefix: d.Cfg.CommandPrefix,
		Middlewares:   mws,
		Token:         d.Cfg.SandboxToken,
	})
	mux.HandleFunc("/sandbox", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/sandbox/", http.StatusFound)
	})

	addr := fmt.Sprintf("%s:%s", d.Cfg.SandboxHost, d.Cfg.SandboxPort)
	srv := &http.Server{Addr: addr, Handler: mux}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	fmt.Printf("sandbox im: open http://%s/sandbox/\n", addr)

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
