package can

import (
	"errors"
	"io"
	"time"
)

// The Reader interface extends the `io.Reader` interface by method
// to read a frame.
type Reader interface {
	io.Reader
	ReadFrame(*Frame) error
}

// The Writer interface extends the `io.Writer` interface by method
// to write a frame.
type Writer interface {
	io.Writer
	WriteFrame(Frame, time.Duration) error
}

// The ReadWriteCloser interface combines the Reader and Writer and
// `io.Closer` interface.
type ReadWriteCloser interface {
	Reader
	Writer

	io.Closer
	setPassFilter(allowedIds []uint32) error
	setBlockFilter(disallowedIds []uint32) error
	deleteFilter() error
}

type readWriteCloser struct {
	rwc    io.ReadWriteCloser
	socket int
}

// NewReadWriteCloser returns a ReadWriteCloser for an `io.ReadWriteCloser`.
func NewReadWriteCloser(rwc io.ReadWriteCloser) ReadWriteCloser {
	return &readWriteCloser{rwc: rwc, socket: 0}
}

func (rwc *readWriteCloser) ReadFrame(frame *Frame) error {
	b := make([]byte, 16)
	n, err := rwc.Read(b)

	if err != nil {
		return err
	}

	err = Unmarshal(b[:n], frame)

	return err
}

func (rwc *readWriteCloser) WriteFrame(frame Frame, min time.Duration) error {
	start := time.Now()
	b, err := Marshal(frame)

	if err != nil {
		return err
	}

	_, err = rwc.Write(b)

	duration := time.Since(start)
	if duration < min {
		time.Sleep(min - duration)
	}

	return err
}

func (rwc *readWriteCloser) Read(b []byte) (n int, err error) {
	return rwc.rwc.Read(b)
}

func (rwc *readWriteCloser) Write(b []byte) (n int, err error) {
	return rwc.rwc.Write(b)
}

func (rwc *readWriteCloser) Close() error {
	rwc.socket = 0
	return rwc.rwc.Close()
}

// ErrorKernelFilterNotSupported is returned if the socket attribute is 0. Then the method
// setsockopt can't be called.
var ErrorKernelFilterNotSupported = errors.New("not possible to set kernel filter")
var ErrorKernelFilterTooMany = errors.New("too many kernel filters")
