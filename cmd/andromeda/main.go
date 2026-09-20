// Command andromeda is the Andromeda-Gateway process entry.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hafuunano/Andromeda-Gateway/internal/config"
	"github.com/Hafuunano/Andromeda-Gateway/internal/driver"
	"github.com/Hafuunano/Andromeda-Gateway/internal/driver/onebot"
	qqopendriver "github.com/Hafuunano/Andromeda-Gateway/internal/driver/qqopen"
	"github.com/Hafuunano/Lucy"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	var drv driver.Driver
	switch cfg.Driver {
	case driver.NameOneBot:
		drv = &onebot.Driver{Cfg: cfg}
	case driver.NameQQOpen:
		drv = &qqopendriver.Driver{Cfg: cfg}
	default:
		log.Fatalf("unsupported DRIVER %q", cfg.Driver)
	}

	result, err := lucy.Bootstrap(lucy.Options{
		NickNames:     cfg.NickNames,
		CommandPrefix: cfg.CommandPrefix,
		SuperUsers:    cfg.SuperUsers,
		PersonaPath:   os.Getenv("SOUL_PERSONA_PATH"),
	})
	if err != nil {
		log.Fatalf("lucy bootstrap: %v", err)
	}

	protocol.Activate(drv.Capabilities())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("andromeda-gateway: starting driver=%s", drv.Name())
	if err := drv.Start(ctx, result.Middlewares); err != nil && ctx.Err() == nil {
		log.Fatalf("driver: %v", err)
	}
}
