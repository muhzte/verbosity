package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var commandDefinitions = []*discordgo.ApplicationCommand{
	{
		Name: "ping",
		Description: "Check if Verbosity is responsive.",
	},
	{
		Name: "join",
		Description: "Make Verbosity join your current voice channel.",
	},
	{
		Name: "leave",
		Description: "Make Verbosity leave its current voice channel.",
	},
}

func (b *Bot) registerCommands() error {
	b.Session.AddHandler(b.handleInteraction)

	for _, cmd := range commandDefinitions {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.cfg.GuildID, cmd)
		if err != nil {
			return err
		}
		log.Printf("Registered command: /%s", cmd.Name)
	}

	return nil
}

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		b.handlePing(s, i)
	case "join":
		b.handleJoin(s, i)
	case "leave":
		b.handleLeave(s, i)
	}
}