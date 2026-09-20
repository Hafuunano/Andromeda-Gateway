// Package driver defines the protocol Driver contract for Andromeda-Gateway.
// Concrete drivers live in subpackages and must declare only real capabilities.
package driver

import (
	"context"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
)

// Name identifies a driver in config (e.g. "onebot", "qqopen").
type Name string

const (
	NameOneBot  Name = "onebot"
	NameQQOpen  Name = "qqopen"
	NameSandbox Name = "sandbox"
)

// Driver starts one protocol ingress and installs the protocol handler chain.
type Driver interface {
	Name() Name
	// Capabilities lists features this driver truly supports. Used by protocol.Activate.
	Capabilities() []protocol.Capability
	// Start blocks or runs the ingress until ctx is cancelled. mws wrap the dispatch chain.
	Start(ctx context.Context, mws []protocol.Middleware) error
}
