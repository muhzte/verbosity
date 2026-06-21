package buffer

import "testing"

func TestRingBufferWrapsAndCaps(t *testing.T) {
	rb := NewRingBuffer()

	for i := 0; i < Capacity+10; i++ {
		rb.Write(Frame{Sequence: uint16(i)})
	}

	got := rb.Snapshot()
	if len(got) != Capacity {
		t.Fatalf("Expected %d frames, got %d", Capacity, len(got))
	}
	if got[0].Sequence != 10 {
		t.Errorf("Oldest frame should be sequence 10, got %d", got[0].Sequence)
	}
	if got[len(got)-1].Sequence != uint16(Capacity+9) {
		t.Errorf("Newest frame should be sequence %d, got %d", Capacity+9, got[len(got)-1].Sequence)
	}
}
