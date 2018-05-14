package nasync

type buffer struct {
	buf []*task // contents are the bytes buf[off : len(buf)]
	off int     // read at &buf[off], write at &buf[len(buf)]
}

func (b *buffer) Reset() {
	b.Truncate(0)
}

func (b *buffer) Len() int { return len(b.buf) - b.off }

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (b *buffer) Truncate(n int) {
	switch {
	case n < 0 || n > b.Len():
		panic("Buffer: truncation out of range")
	case n == 0:
		// Reuse buffer space.
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+n]
}

func (b *buffer) Tasks() []*task {
	return b.buf[b.off:]
}

func (b *buffer) Append(t *task) {
	b.buf = append(b.buf, t)
}
