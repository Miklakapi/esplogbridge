package bridge

import (
	"context"
	"fmt"
	"time"

	"github.com/Miklakapi/esplogbridge/internal/config"
)

type Bridge struct {
	cfg   config.Config
	queue chan Event
	loki  *Loki
}

func New(cfg config.Config) *Bridge {
	return &Bridge{
		cfg:   cfg,
		queue: make(chan Event, cfg.Input.UDP.QueueSize),
		loki:  NewLoki(cfg),
	}
}

func (b *Bridge) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- runUDPReceiver(ctx, b.cfg, b.queue)
	}()

	maxItems := b.cfg.Loki.Batch.MaxItems
	maxWait := b.cfg.Loki.Batch.MaxWait

	t := time.NewTimer(maxWait)
	defer t.Stop()

	batch := make([]Event, 0, maxItems)

	for {
		select {
		case <-ctx.Done():
			b.flushBatch(ctx, batch)
			return nil

		case err := <-errCh:
			b.flushBatch(ctx, batch)
			if err == nil || ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("udp receiver: %w", err)

		case ev := <-b.queue:
			batch = append(batch, ev)
			if len(batch) >= maxItems {
				b.flushBatch(ctx, batch)
				batch = batch[:0]
				if !t.Stop() {
					select {
					case <-t.C:
					default:
					}
				}
				t.Reset(maxWait)
			}

		case <-t.C:
			b.flushBatch(ctx, batch)
			batch = batch[:0]
			t.Reset(maxWait)
		}
	}
}

func (b *Bridge) flushBatch(ctx context.Context, batch []Event) {
	if len(batch) == 0 {
		return
	}
	_ = b.loki.SendBatch(ctx, batch)
}
