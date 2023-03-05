package queue

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	TimeoutError = errors.New("timeout")
)

type MessageQueue[T any] struct {
	sync.Mutex
	messages List[T]
	readers  List[chan T]
}

func New[T comparable]() MessageQueue[T] {
	return MessageQueue[T]{}
}

func (mq *MessageQueue[T]) Push(v T) {
	mq.Lock()
	defer mq.Unlock()

	if reader := mq.readers.Pop(); reader != nil {
		reader.value <- v
		return
	}
	mq.messages.Push(v)
}

func (mq *MessageQueue[T]) Pop(timeout time.Duration) (T, error) {
	mq.Lock()

	if item := mq.messages.Pop(); item != nil {
		mq.Unlock()
		return item.value, nil
	}
	var (
		ctx    = context.Background()
		cancel context.CancelFunc
		ch     = make(chan T)
		v      T
	)
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	i := mq.readers.Push(ch)
	mq.Unlock()
	select {
	case <-ctx.Done():
		mq.readers.Remove(i)
		close(ch)
		return v, TimeoutError
	case v = <-ch:
		return v, nil
	}
}
