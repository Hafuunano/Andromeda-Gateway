package onebot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Hafuunano/Andromeda-Gateway/internal/config"
	"github.com/Hafuunano/Andromeda-Gateway/internal/driver"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol/zerobot"
	zero "github.com/wdvxdr1123/ZeroBot"
	zbdriver "github.com/wdvxdr1123/ZeroBot/driver"
)

// Driver runs ZeroBot WebSocket ingress and installs the zerobot protocol adaptor.
type Driver struct {
	Cfg config.Config
}

// Name implements driver.Driver.
func (d *Driver) Name() driver.Name { return driver.NameOneBot }

// Capabilities implements driver.Driver.
func (d *Driver) Capabilities() []protocol.Capability {
	return protocol.DefaultOneBotCapabilities()
}

// Start implements driver.Driver. Blocks until ZeroBot stops (typically forever).
func (d *Driver) Start(ctx context.Context, mws []protocol.Middleware) error {
	_ = ctx
	if d.Cfg.WSURL == "" {
		return fmt.Errorf("onebot: WS_URL is required")
	}
	zerobot.InstallWithMiddlewares(mws)

	superUsers := make([]int64, 0, len(d.Cfg.SuperUsers))
	for _, s := range d.Cfg.SuperUsers {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil || n == 0 {
			continue
		}
		superUsers = append(superUsers, n)
	}

	zero.RunAndBlock(&zero.Config{
		NickName:      append([]string(nil), d.Cfg.NickNames...),
		CommandPrefix: d.Cfg.CommandPrefix,
		SuperUsers:    superUsers,
		Driver: []zero.Driver{
			zbdriver.NewWebSocketServer(16, d.Cfg.WSURL, d.Cfg.WSToken),
		},
	}, nil)
	return nil
}
