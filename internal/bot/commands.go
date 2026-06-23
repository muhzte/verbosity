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
	discord.SlashCommandCreate{
		Name:        "bufferstatus",
		Description: "Show how much audio Verbosity has buffered for you.",
	},
	discord.SlashCommandCreate{
		Name: "clip",
		Description: "Export a voice evidence clip for a user.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name: "user",
				Description: "The user to clip.",
				Required: true,
			},
			discord.ApplicationCommandOptionInt{
				Name: "seconds",
				Description: "How many seconds to clip (1-30). Defaults to 30.",
				Required: false,
			},
		},
	},
}

func registerCommands(ctx context.Context, client *bot.Client, guildID snowflake.ID) error {
	_, err := client.Rest.SetGuildCommands(client.ApplicationID, guildID, commandDefinitions)
	return err
}
