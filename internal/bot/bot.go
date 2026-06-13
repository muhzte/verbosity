package bot

import (
	"log"
	"github.com/bwmarrin/discordgo"
	"github.com/muhzte/verbosity/internal/config"
)

type Bot struct {
	Session *discordgo.Session
	cfg *config.Config
}

func New(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		Session: session,
		cfg: cfg,
	}

	session.AddHandler(b.onReady)

	return b, nil
}

func (b *Bot) onReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Verbosity is online! Logged in as %s#%s", r.User.Username, r.User.Discriminator)
}

func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return err
	}

	return b.registerCommands()
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}