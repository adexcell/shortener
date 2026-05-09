package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/adexcell/shortener/internal/shortener/dto"
)

// Redis defines the interface for setting original URL in the cache.
type Redis interface {
	Set(ctx context.Context, key, value string) error
}

// AsyncRedisWriter handles asynchronous updates to the original URL cache.
type AsyncRedisWriter struct {
	redis         Redis
	ch            chan dto.ShortenTask
	isClosed      atomic.Bool
	wg            sync.WaitGroup
	drainTimeout time.Duration
}

// NewAsyncRedisWriter creates a new instance of AsyncRedisWriter.
func NewAsyncRedisWriter(redis Redis) *AsyncRedisWriter {
	return &AsyncRedisWriter{
		redis:         redis,
		ch:            make(chan dto.ShortenTask, 1000),
		drainTimeout: 10 * time.Second,
	}
}

// Start begins the background worker for processing the Redis update queue.
func (w *AsyncRedisWriter) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.worker(ctx)
}

// Send queues a status update task to be processed asynchronously.
func (w *AsyncRedisWriter) Send(ctx context.Context, task dto.ShortenTask) {
	select {
	case <-ctx.Done():
		log.Debug().Msgf("AsyncRedisWriter send cancelled with context.Done(): %s", ctx.Err().Error())
		return
	default:
	}
	if !w.isClosed.Load() {
		w.ch <- task
	}
}

// Stop gracefully shuts down the writer by closing the queue and waiting for workers to finish.
func (w *AsyncRedisWriter) Stop() {
	if w.isClosed.CompareAndSwap(false, true) {
		close(w.ch)
	}
	w.wg.Wait()
}

func (w *AsyncRedisWriter) worker(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			w.drainRemaining()
			return
		case task, ok := <-w.ch:
			if !ok {
				return // channel closed
			}
			w.handleTask(ctx, task)
		}
	}
}

func (w *AsyncRedisWriter) handleTask(ctx context.Context, task dto.ShortenTask) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := w.redis.Set(ctx, task.Shorten, task.OriginalURL)
	if err != nil {
		log.Debug().Err(err).Str("id", task.Shorten).Msg("[async-redis-writer] failed to set")
	}

}

func (w *AsyncRedisWriter) drainRemaining() {
	timer := time.NewTimer(w.drainTimeout)
	defer timer.Stop()

	ctx := context.Background()

	for {
		select {
		case task, ok := <-w.ch:
			if !ok {
				return
			}
			w.handleTask(ctx, task)
		case <-timer.C:
			return
		}
	}
}
