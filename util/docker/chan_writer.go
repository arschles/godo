package docker

// ChanWriter returns an io.Writer that sends all of its writes as log messages to ch
type ChanWriter struct {
	ch chan<- Log
}

// NewChanWriter creates a new ChanWriter ready to accept writes and forward
func NewChanWriter(ch chan<- Log) *ChanWriter {
	return &ChanWriter{ch: ch}
}

// Write is the io.Writer interface implementation
func (c *ChanWriter) Write(b []byte) (int, error) {
	c.ch <- Log{str: string(b)}
	return len(b), nil
}
