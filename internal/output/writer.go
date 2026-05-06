package output

import (
	"fmt"
	"io"
	"os"
)

// Destination represents where output should be written.
type Destination string

const (
	DestStdout Destination = "stdout"
	DestStderr Destination = "stderr"
	DestFile   Destination = "file"
)

// Writer wraps an io.Writer with destination metadata.
type Writer struct {
	dest Destination
	w    io.Writer
}

// NewWriter creates a Writer for the given destination.
// For DestFile, path must be non-empty.
func NewWriter(dest Destination, path string) (*Writer, error) {
	switch dest {
	case DestStdout:
		return &Writer{dest: dest, w: os.Stdout}, nil
	case DestStderr:
		return &Writer{dest: dest, w: os.Stderr}, nil
	case DestFile:
		if path == "" {
			return nil, fmt.Errorf("file destination requires a non-empty path")
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, fmt.Errorf("opening output file %q: %w", path, err)
		}
		return &Writer{dest: dest, w: f}, nil
	default:
		return nil, fmt.Errorf("unknown destination %q", dest)
	}
}

// Write writes p to the underlying writer.
func (w *Writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

// Destination returns the destination type of this writer.
func (w *Writer) Destination() Destination {
	return w.dest
}

// Close closes the underlying writer if it implements io.Closer.
func (w *Writer) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
