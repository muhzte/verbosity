package audio

import (
	// imports in each file are alphabetical by the way
	"fmt"
	"os"
	"time"

	"github.com/muhzte/verbosity/internal/buffer"
)

type Clip struct {
	Path string
	Duration float64
	FrameCount int
	CapturedAt time.Time
}

func Export(frames []buffer.Frame) (*Clip, error) {
	if len(frames) == 0 {
		return nil, fmt.Errorf("No frames to export")
	}

	f, err := os.CreateTemp("", "verbosity-clip-*.ogg")
	if err != nil {
		return nil, fmt.Errorf("Create temp file: %w", err)
	}
	defer f.Close()

	raw := make([][]byte, len(frames))
	for i, fr := range frames {
		raw[i] = fr.Opus
	}

	if err := writeOpusOGG(f, raw); err != nil {
		os.Remove(f.Name())
		return nil, fmt.Errorf("Write ogg: %w", err)
	}

	return &Clip{
		Path: f.Name(),
		Duration: float64(len(frames)) / float64(buffer.FrameRate),
		FrameCount: len(frames),
		CapturedAt: time.Now().UTC(),
	}, nil
}