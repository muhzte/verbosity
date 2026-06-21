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
	voiceState, ok := b.Client.Caches.VoiceState(*e.GuildID(), e.User().ID)
	if !ok || voiceState.ChannelID == nil {
		respond(e, "You need to be in a voice channel for me to join.")
		return
	}

	conn := b.Client.VoiceManager.CreateConn(*e.GuildID())
	conn.SetOpusFrameReceiver(verbosityvoice.NewBufferReceiver(b.bufferMgr))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.Open(ctx, *voiceState.ChannelID, false, false); err != nil {
		log.Printf("Join: Failed to open voice connection: %v", err)
		respond(e, "I couldn't join your voice channel.")
		return
	}

	respond(e, "Joined your voice channel.")
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
func respond(e *events.ApplicationCommandInteractionCreate, message string) {
	if err := e.CreateMessage(discord.MessageCreate{Content: message}); err != nil {
		log.Printf("Respond: Failed to send: %v", err)
	}
}
