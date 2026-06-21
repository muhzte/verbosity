package buffer

import "sync"

const (
	FrameRate     = 50
	BufferSeconds = 30
	Capacity      = FrameRate * BufferSeconds
)

type Frame struct {
	Opus      []byte
	Sequence  uint16
	Timestamp uint32
}

type RingBuffer struct {
	mu     sync.Mutex
	frames []Frame
	next   int
	filled int
}

func NewRingBuffer() *RingBuffer {
	return &RingBuffer{frames: make([]Frame, Capacity)}
}

func (r *RingBuffer) Write(f Frame) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.frames[r.next] = f
	r.next = (r.next + 1) % Capacity
	if r.filled < Capacity {
		r.filled++
	}
}

func (r *RingBuffer) Snapshot() []Frame {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Frame, r.filled)
	if r.filled < Capacity {
		copy(out, r.frames[:r.filled])
		return out
	}

	copy(out, r.frames[r.next:])
	copy(out[Capacity-r.next:], r.frames[:r.next])
	return out
}
