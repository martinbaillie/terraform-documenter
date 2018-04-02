package markdown

import (
	"bytes"
	"io"
)

type prefixWriter struct {
	w      io.Writer
	prefix []byte

	wrotePrefix bool
}

// newPrefixWriter creates a new prefix writer that prefixes every line with a prefix.
func newPrefixWriter(w io.Writer, prefix []byte) io.Writer {
	return &prefixWriter{
		w:      w,
		prefix: prefix,
	}
}

func (iw *prefixWriter) Write(p []byte) (n int, err error) {
	for i, b := range p {
		if !iw.wrotePrefix {
			prefix := iw.prefix
			if b == '\n' {
				prefix = bytes.TrimSpace(iw.prefix)
			}
			_, err = iw.w.Write(prefix)
			if err != nil {
				return n, err
			}
			n += len(prefix)
			iw.wrotePrefix = true
		}
		_, err = iw.w.Write(p[i : i+1])
		if err != nil {
			return n, err
		}
		if b == '\n' {
			iw.wrotePrefix = false
		}
		n++
	}
	return len(p), nil
}
