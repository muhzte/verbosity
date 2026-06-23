package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"

	"github.com/muhzte/verbosity/internal/audio"
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
	case "clip":
		b.HandleClip(e)
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
	
	if err := e.DeferCreateMessage(false); err != nil {
		log.Printf("Join: Failed to defer response: %v", err)
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

func (b *Bot) handleClip(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()

	targetID := data.Snowflake("user")

	seconds := 30
	if s, ok := data.OptInt("seconds"); ok && s > 0 && s <= 30 {
		seconds = s
	}

	if err := e.DeferCreateMessage(false); err != nil {
		log.Printf("Clip: Defer failed: %v", err)
		return
	}

	frames := b.bufferMgr.Snapshot(targetID)
	if len(frames) == 0 {
		followUp(e, fmt.Sprintf("No audio buffered for <@%d> - are they in a monitored voice channel?", targetID))
		return
	}

	maxFrames := seconds * buffer.FrameRate
	if len(frames) > maxFrames {
		frames = frames[len(frames)-maxFrames:]
	}

	clip, err := audio.Export(frames)
	if err != nil {
		log.Printf("Clip: Export Failed: %v", err)
		followUp(e, "Failed to export audio clip.")
		return
	}
	defer os.Remove(clip.Path)

	f, err := os.Open(clip.Path)
	if err != nil {
		log.Printf("Clip: Open temp file: %v", err)
		followUp(e, "Failed to read exported clip.")
		return
	}
	defer f.Close()

	displayName := fmt.Sprintf("<@%d>", targetID)
	if member, ok := e.Client().Caches.Member(*e.GuildID(), targetID); ok {
		displayName = member.User.Username
	}

	msg := fmt.Sprintf(
		"📎 **Voice Clip** - %s\n**Duration:** %.1fs (%d frames)\n**Captured:** %s UTC",
		displayName,
		clip.Duration,
		clip.FrameCount,
		clip.CapturedAt.Format(time.DateTime),
	)

	filename := fmt.Sprintf("clip-%d-%d.ogg", targetID, clip.CapturedAt.Unix())

	_, err = e.Client().Rest.CreateFollowupMessage(
		e.ApplicationID(),
		e.Token(),
		discord.NewMessageCreateBuilder().
		SetContent(msg).
		AddFile(filename, f).
		Build(),
	)
	if err != nil {
		log.Printf("Clip: Send Failed: %v", err)
	}
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
