package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"

	"github.com/muhzte/verbosity/internal/buffer"
	verbosityvoice "github.com/muhzte/verbosity/internal/voice"
)

func (b *Bot) handleCommand(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()

	switch data.CommandName() {
	case "ping":
		b.handlePing(e)
	case "join":
		b.handleJoin(e)
	case "leave":
		b.handleLeave(e)
	case "bufferstatus":
		b.handleBufferStatus(e)
	}
}

func (b *Bot) handlePing(e *events.ApplicationCommandInteractionCreate) {
	if err := e.CreateMessage(discord.MessageCreate{Content: "Pong! Verbosity is online."}); err != nil {
		log.Printf("Ping: Failed to Respond: %v", err)
	}
}

func (b *Bot) handleJoin(e *events.ApplicationCommandInteractionCreate) {
	if err := e.DeferCreateMessage(false); err != nil {
		log.Printf("Join: Failed to defer response: %v", err)
		return
	}

	voiceState, ok := b.Client.Caches.VoiceState(*e.GuildID(), e.User().ID)
	if !ok || voiceState.ChannelID == nil {
		respond(e, "You need to be in a voice channel for me to join.")
		return
	}

	channelID := *voiceState.ChannelID
	guildID := *e.GuildID()

	go func() {
		conn := b.Client.VoiceManager.CreateConn(guildID)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		if err := conn.Open(ctx, channelID, false, false); err != nil {
			log.Printf("Join: Voice open failed: %v", err)
			followUp(e, "I couldn't join your voice channel.")
			return
		}

		conn.SetOpusFrameReceiver(verbosityvoice.NewBufferReceiver(b.bufferMgr))
		followUp(e, "Joined your voice channel.")
	}()
}

func (b *Bot) handleLeave(e *events.ApplicationCommandInteractionCreate) {
	conn := b.Client.VoiceManager.GetConn(*e.GuildID())
	if conn == nil {
		respond(e, "I'm not currently in a voice channel here.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn.Close(ctx)

	respond(e, "Left the voice channel.")
}

func (b *Bot) handleBufferStatus(e *events.ApplicationCommandInteractionCreate) {
	frames := b.bufferMgr.Snapshot(e.User().ID)
	seconds := float64(len(frames)) / float64(buffer.FrameRate)
	respond(e, fmt.Sprintf("%d frames buffered (%.1fs) for you.", len(frames), seconds))
}
func followUp(e *events.ApplicationCommandInteractionCreate, message string) {
	_, err := e.Client().Rest.CreateFollowupMessage(
		e.ApplicationID(),
		e.Token(),
		discord.MessageCreate{Content: message},
	)
	if err != nil {
		log.Printf("followup: Failed to send: %v", err)
	}
}

func respond(e *events.ApplicationCommandInteractionCreate, message string) {
	if err := e.CreateMessage(discord.MessageCreate{Content: message}); err != nil {
		log.Printf("Respond: Failed to send: %v", err)
	}
}
