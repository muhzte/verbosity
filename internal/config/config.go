package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	GuildID string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN is not set (check your .env file)")
	}
	guildID := os.Getenv("DISCORD_GUILD_ID")
	if guildID == "" {
		return nil, fmt.Errorf("DISCORD_GUILD_ID is not set (ckeck the .env file)")
	}

	return &Config{
		DiscordToken: token,
		GuildID: guildID,
	}, nil
}