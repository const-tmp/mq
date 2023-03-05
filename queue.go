package queue

import (
	"context"
	"errors"
	"sync"
	"time"
)

var TimeoutError = errors.New("timeout")

type MessageQueue[T any] struct {
	sync.Mutex
	messages fifo[T]
	readers  fifo[chan T]
}

func (mq *MessageQueue[T]) Push(v T) {
	mq.Lock()
	defer mq.Unlock()
	if reader := mq.readers.pop(); reader != nil {
		reader.value <- v
		return
	}
	mq.messages.push(v)
}

func (mq *MessageQueue[T]) Pop(timeout time.Duration) (T, error) {
	mq.Lock()
	if item := mq.messages.pop(); item != nil {
		mq.Unlock()
		return item.value, nil
	}
	var (
		ctx    = context.Background()
		cancel context.CancelFunc
		v      T
	)
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	ch := make(chan T)
	i := mq.readers.push(ch)
	mq.Unlock()
	select {
	case <-ctx.Done():
		i.expired = true
		close(ch)
		return v, TimeoutError
	case v = <-ch:
		return v, nil
	}
}
