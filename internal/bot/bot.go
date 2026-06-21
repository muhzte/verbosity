package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/godave/golibdave"
	"github.com/disgoorg/snowflake/v2"
	"github.com/muhzte/verbosity/internal/config"
)

type Bot struct {
	Client *bot.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Bot, error) {
	b := &Bot{cfg: cfg}

	client, err := disgo.New(cfg.DiscordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildVoiceStates,
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagsAll),
		),
		bot.WithVoiceManagerConfigOpts(

			voice.WithDaveSessionCreateFunc(golibdave.NewSession),
		),
		bot.WithEventListenerFunc(b.onReady),
		bot.WithEventListenerFunc(b.handleCommand),
	)
	if err != nil {
		return nil, err
	}

	b.Client = client
	return b, nil
}

func (b *Bot) onReady(e *events.Ready) {
	log.Printf("Verbosity is online! Logged in as %s", e.User.Username)
}

func (b *Bot) Start(ctx context.Context) error {
	guildID, err := snowflake.Parse(b.cfg.GuildID)
	if err != nil {
		return fmt.Errorf("Invalid DISCORD_GUILD_ID: %w", err)
	}

	if err := registerCommands(ctx, b.Client, guildID); err != nil {
		return err
	}
	return b.Client.OpenGateway(ctx)
}

func (b *Bot) Stop(ctx context.Context) {
	b.Client.Close(ctx)
}
