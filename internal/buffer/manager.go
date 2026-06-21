package buffer

import (
	"sync"

	"github.com/disgoorg/snowflake/v2"
)

type Manager struct {
	mu      sync.Mutex
	buffers map[snowflake.ID]*RingBuffer
}

func NewManager() *Manager {
	return &Manager{buffers: make(map[snowflake.ID]*RingBuffer)}
}

func (m *Manager) Write(userID snowflake.ID, f Frame) {
	m.getOrCreate(userID).Write(f)
}

func (m *Manager) Snapshot(userID snowflake.ID) []Frame {
	m.mu.Lock()
	buf, ok := m.buffers[userID]
	m.mu.Unlock()

	if !ok {
		return nil
	}
	return buf.Snapshot()
}

func (m *Manager) getOrCreate(userID snowflake.ID) *RingBuffer {
	m.mu.Lock()
	defer m.mu.Unlock()

	buf, ok := m.buffers[userID]
	if !ok {
		buf = NewRingBuffer()
		m.buffers[userID] = buf
	}
	return buf
}
