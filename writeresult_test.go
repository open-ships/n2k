package n2k

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteResult_WaitContext(t *testing.T) {
	r := newWriteResult()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.ErrorIs(t, r.WaitContext(ctx), context.Canceled)

	r.complete(nil)
	assert.NoError(t, r.WaitContext(context.Background()))
}

func TestWriteResult_Wait(t *testing.T) {
	r := newWriteResult()
	go func() {
		time.Sleep(10 * time.Millisecond)
		r.complete(nil)
	}()
	err := r.Wait()
	assert.NoError(t, err)
}

func TestWriteResult_WaitWithError(t *testing.T) {
	r := newWriteResult()
	go func() {
		r.complete(fmt.Errorf("write failed"))
	}()
	err := r.Wait()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write failed")
}

func TestWriteResult_Done(t *testing.T) {
	r := newWriteResult()
	go func() {
		time.Sleep(10 * time.Millisecond)
		r.complete(nil)
	}()

	select {
	case <-r.Done():
		// success
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Done")
	}
	assert.NoError(t, r.Wait())
}

func TestWriteResult_ZeroValue(t *testing.T) {
	var r WriteResult
	assert.NoError(t, r.Wait())

	select {
	case <-r.Done():
		// should be immediately ready
	default:
		t.Fatal("zero-value Done() should be immediately ready")
	}
}

func TestWriteResult_MultipleWaits(t *testing.T) {
	r := newWriteResult()
	r.complete(nil)
	// Multiple calls to Wait should all succeed
	assert.NoError(t, r.Wait())
	assert.NoError(t, r.Wait())
}
