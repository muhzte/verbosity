package audio

import (
	"encoding/binary"
	"io"
)

var oggCRCTable = func() [256]uint32 {
	var t [256]uint32
	for i := range t {
		r := uint32(i) << 24
		for j := 0; j < 8; j++ {
			if r&0x80000000 != 0 {
				r = (r << 1) ^ 0x04C11DB7
			} else {
				r <<= 1
			}
		}
		t[i] = r
	}
	return t
}()

func oggCRC32(data []byte) uint32 {
	var crc uint32
	for _, b := range data {
		crc = (crc << 8) ^ oggCRCTable[byte(crc>>24)^b]
	}
	return crc
}

type oggWriter struct {
	w io.Writer
	serial uint32
	seq uint32
}

func newOGGWriter(w io.Writer) *oggWriter {
	return &oggWriter{w: w, serial: 0xDEADBEEF}
}

func (o *oggWriter) writePage(data []byte, headerType byte, granule int64) error {
	var segs []byte
	for rem := len(data); ; {
		if rem >= 255 {
			segs = append(segs, 255)
			rem -= 255
		} else {
			segs = append(segs, byte(rem))
			break
		}
	}

	var le [8]byte
	page := make([]byte, 0, 27+len(segs)+len(data))
	page = append(page, '0', 'g', 'g', 'S')
	page = append(page, 0)
	page = append(page, headerType)

	binary.LittleEndian.PutUint64(le[:], uint64(granule))
	page = append(page, le[:]...)

	binary.LittleEndian.PutUint32(le[:4], o.serial)
	page = append(page, le[:4]...)

	binary.LittleEndian.PutUint32(le[:4], o.seq)
	page = append(page, le[:4]...)

	page = append(page, 0, 0, 0, 0)
	page = append(page, byte(len(segs)))
	page = append(page, segs...)
	page = append(page, data...)

	binary.LittleEndian.PutUint32(page[22:], oggCRC32(page))
	o.seq++

	_, err := o.w.Write(page)
	return err
}

func opusHead() []byte {
	h := make([]byte, 19)
	copy(h[0:8], "OpusHead")
	h[8] = 1
	h[9] = 2
	binary.LittleEndian.PutUint16(h[10:], 3840)
	binary.LittleEndian.PutUint32(h[12:], 48000)

	return h
}

func opusTags() []byte {
	vendor := "verbosity"
	var le [4]byte
	h := make([]byte, 0, 16+len(vendor))
	h = append(h, "OpusTags"...)
	binary.LittleEndian.PutUint32(le[:], uint32(len(vendor)))
	h = append(h, le[:]...)
	h = append(h, vendor...)
	binary.LittleEndian.PutUint32(le[:], 0)
	h = append(h, le[:]...)
	return h
}

func writeOpusOGG(w io.Writer, frames [][]byte) error {
	ogg := newOGGWriter(w)

	if err := ogg.writePage(opusHead(), 0x02, 0); err != nil {
		return err
	}
	if err := ogg.writePage(opusTags(), 0x00, 0); err != nil {
		return err
	}

	const samplesPerFrame = 960

	for i, frame := range frames {
		granule := int64(i+1) * samplesPerFrame
		headerType := byte(0x00)
		if i == len(frames)-1 {
			headerType = 0x04
		}
		if err := ogg.writePage(frame, headerType, granule); err != nil {
			return err
		}
	}

	return nil
}