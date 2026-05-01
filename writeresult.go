package n2k

// WriteResult represents the outcome of an asynchronous write operation.
// The zero value is safe to use and behaves as a completed, successful write.
type WriteResult struct {
	done chan struct{}
	err  error
}

// newWriteResult creates a WriteResult whose done channel is open (not yet complete).
func newWriteResult() *WriteResult {
	return &WriteResult{
		done: make(chan struct{}),
	}
}

// complete marks the write as finished and records any error.
func (r *WriteResult) complete(err error) {
	r.err = err
	close(r.done)
}

// Wait blocks until the write is complete, then returns the error (or nil).
// On a zero-value WriteResult, Wait returns nil immediately.
func (r *WriteResult) Wait() error {
	if r.done == nil {
		return nil
	}
	<-r.done
	return r.err
}

// Done returns a channel that is closed when the write is complete.
// On a zero-value WriteResult, Done returns an already-closed channel.
func (r *WriteResult) Done() <-chan struct{} {
	if r.done == nil {
		ch := make(chan struct{})
		close(ch)
		return ch
	}
	return r.done
}
