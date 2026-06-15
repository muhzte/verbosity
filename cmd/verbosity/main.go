package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/muhzte/verbosity/internal/bot"

	"github.com/muhzte/verbosity/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	ctx := context.Background()
	if err := b.Start(ctx); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	defer b.Stop(ctx)

	log.Println("Verbosity is running. Press Ctrl+C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down.")
}