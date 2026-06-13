package bot

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respond(s, i, "Pong. Verbosity is online.")
}

func (b *Bot) handleJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.State.Guild(i.GuildID)
	if err != nil {
		respond(s, i, "I couldn't look up this server's info.")
		return
	}

	var voiceChannelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	if voiceChannelID == "" {
		respond(s, i, "You need to be in a voice channel for me to join.")
		return
	}

	_, err = s.ChannelVoiceJoin(i.GuildID, voiceChannelID, false, false)
	if err != nil {
		respond(s, i, "I couldn't join your voice channel.")
		return
	}

	respond(s, i, "Joined your voice channel.")
}

func (b *Bot) handleLeave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	vc, ok := s.VoiceConnections[i.GuildID]
	if !ok {
		respond(s, i, "I'm not currently in a voice channel here.")
		return
	}

	if err := vc.Disconnect(); err != nil {
		respond(s, i, "I had trouble leaving the voice channel.")
		return
	}

	respond(s, i, "Left the voice channel.")
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}