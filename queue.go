package mq

import (
	"context"
	"errors"
	"sync"
	"time"
)

var TimeoutError = errors.New("timeout")

type MQ[T any] struct {
	sync.Mutex
	messages fifo[T]
	readers  fifo[chan T]
}

func (mq *MQ[T]) Push(v T) {
	mq.Lock()
	defer mq.Unlock()
	if reader := mq.readers.pop(); reader != nil {
		reader.value <- v
		return
	}
	mq.messages.push(v)
}

func (mq *MQ[T]) Pop(timeout time.Duration) (T, error) {
	mq.Lock()
	if item := mq.messages.pop(); item != nil {
		mq.Unlock()
		return item.value, nil
	}
	ch := make(chan T)
	i := mq.readers.push(ch)
	mq.Unlock()
	var (
		ctx    = context.Background()
		cancel context.CancelFunc
		v      T
	)
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	select {
	case <-ctx.Done():
		mq.Lock()
		defer mq.Unlock()
		i.expired = true
		close(ch)
		return v, TimeoutError
	case v = <-ch:
		return v, nil
	}
}
