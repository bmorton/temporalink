package temporalink

import (
	"context"
	"testing"
	"time"
)

func TestNewEmbeddedTemporal(t *testing.T) {
	_, err := NewEmbeddedTemporal("0.0.0.0", 7233, 8088)
	if err != nil {
		t.Error(err)
	}
}

func TestEmbeddedTemporal_Start(t *testing.T) {
	et, err := NewEmbeddedTemporal("127.0.0.1", 7233, 8088)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.NewTicker(5 * time.Second).C
		cancel()
	}()
	err = et.Start(ctx)
	if err != nil {
		t.Error(err)
	}
}
