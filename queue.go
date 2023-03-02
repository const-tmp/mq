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

type queue[T comparable] struct {
	sync.Mutex
	slice []T
}

func (q *queue[T]) push(v T) {
	q.Lock()
	q.slice = append(q.slice, v)
	q.Unlock()
}

func (q *queue[T]) pop() (v T, ok bool) {
	q.Lock()
	defer q.Unlock()
	if len(q.slice) == 0 {
		return
	}
	v, q.slice = q.slice[0], q.slice[1:]
	ok = true
	return
}

func (q *queue[T]) remove(v T) {
	q.Lock()
	defer q.Unlock()
	for i, t := range q.slice {
		if t == v {
			q.slice = append(q.slice[:i], q.slice[i+1:]...)
			return
		}
	}
}

type MessageQueue[T comparable] struct {
	messages *queue[T]
	readers  *queue[chan T]
}

func New[T comparable]() MessageQueue[T] {
	return MessageQueue[T]{messages: new(queue[T]), readers: new(queue[chan T])}
}

func (mq MessageQueue[T]) Push(v T) {
	if reader, ok := mq.readers.pop(); ok {
		reader <- v
		return
	}
	mq.messages.push(v)
}

func (mq MessageQueue[T]) Pop(timeout time.Duration) (*T, error) {
	if v, ok := mq.messages.pop(); ok {
		return &v, nil
	}
	var (
		ctx    = context.Background()
		cancel context.CancelFunc
	)
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	ch := make(chan T)
	mq.readers.push(ch)
	select {
	case <-ctx.Done():
		mq.readers.remove(ch)
		close(ch)
		return nil, TimeoutError
	case v := <-ch:
		return &v, nil
	}
}
