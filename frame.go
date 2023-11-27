package can

import (
	"bytes"
	"encoding/binary"
)

const (
	// MaxFrameDataLength defines the max length of a CAN data frame defined in ISO 11898-1.
	MaxFrameDataLength = 8
	// MaxExtFrameDataLength defines the max length of an CAN extended data frame defined in ISO ISO 11898-7.
	MaxExtFrameDataLength = 64
)

// Frame represents a standard CAN data frame
type Frame struct {
	// bit 0-28: CAN identifier (11/29 bit)
	// bit 29: error message flag (ERR)
	// bit 30: remote transmision request (RTR)
	// bit 31: extended frame format (EFF)
	ID     uint32
	Length uint8
	Flags  uint8
	Res0   uint8
	Res1   uint8
	Data   [MaxFrameDataLength]uint8
}

func (from *Frame) Equals(to *Frame) bool {
	if from.ID == to.ID {
		if from.Length == to.Length {
			for i := 0; i < int(from.Length); i++ {
				if from.Data[i] != to.Data[i] {
					return false
				}
			}
			return true
		}
	}

	return false
}

// Marshal returns the byte encoding of frm.
func Marshal(frm Frame) (b []byte, err error) {
	wr := errWriter{
		buf: bytes.NewBuffer([]byte{}),
	}
	wr.write(&frm.ID)
	wr.write(&frm.Length)
	wr.write(&frm.Flags)
	wr.write(&frm.Res0)
	wr.write(&frm.Res1)
	wr.write(&frm.Data)

	return wr.buf.Bytes(), wr.err
}

// Unmarshal parses the bytes b and stores the result in the value
// pointed to by frm.
func Unmarshal(b []byte, frm *Frame) (err error) {
	cr := &errReader{
		buf: bytes.NewBuffer(b),
	}
	cr.read(&frm.ID)
	cr.read(&frm.Length)
	cr.read(&frm.Flags)
	cr.read(&frm.Res0)
	cr.read(&frm.Res1)
	cr.read(&frm.Data)

	return cr.err
}

type errReader struct {
	buf *bytes.Buffer
	err error
}

func (r *errReader) read(v interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.buf, binary.LittleEndian, v)
	}
}

type errWriter struct {
	buf *bytes.Buffer
	err error
}

func (wr *errWriter) write(v interface{}) {
	if wr.err == nil {
		wr.err = binary.Write(wr.buf, binary.LittleEndian, v)
	}
}
