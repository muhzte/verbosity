package voice

import (
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"

	"github.com/muhzte/verbosity/internal/buffer"
)

type BufferReceiver struct {
	bufferMgr *buffer.Manager
}

func NewBufferReceiver(bufferMgr *buffer.Manager) *BufferReceiver {
	return &BufferReceiver{bufferMgr: bufferMgr}
}

func (r *BufferReceiver) ReceiveOpusFrame(userID snowflake.ID, packet *voice.Packet) error {
	opus := make([]byte, len(packet.Opus))
	copy(opus, packet.Opus)

	r.bufferMgr.Write(userID, buffer.Frame{
		Opus:      opus,
		Sequence:  packet.Sequence,
		Timestamp: packet.Timestamp,
	})
	return nil
}

func (r *BufferReceiver) CleanupUser(userID snowflake.ID) {}
func (r *BufferReceiver) Close()
