package bot

import (
	"context"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var commandDefinitions = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "ping",
		Description: "Check if Verbosity is responsive.",
	},
	discord.SlashCommandCreate{
		Name:        "join",
		Description: "Make Verbosity join your current voice channel.",
	},
	discord.SlashCommandCreate{
		Name:        "leave",
		Description: "Make Verbosity leave its current voice channel.",
	},
}

func registerCommands(ctx context.Context, client *bot.Client, guildID snowflake.ID) error {
	_, err := client.Rest.SetGuildCommands(client.ApplicationID, guildID, commandDefinitions)
	return err
}
