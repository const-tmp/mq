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

type Queue[T any] struct {
	sync.Mutex
	list *List[T]
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{list: new(List[T])}
}

func (q *Queue[T]) push(v T) *ListItem[T] {
	q.Lock()
	defer q.Unlock()
	return q.list.Push(v)
}

func (q *Queue[T]) pop() (v T, ok bool) {
	q.Lock()
	defer q.Unlock()
	if i := q.list.Pop(); i == nil {
		return
	} else {
		v, ok = i.value, true
		return
	}
}

func (q *Queue[T]) remove(i *ListItem[T]) {
	q.Lock()
	defer q.Unlock()
	i = i.next
}

type MessageQueue[T any] struct {
	messages *Queue[T]
	readers  *Queue[chan T]
}

func New[T comparable]() MessageQueue[T] {
	return MessageQueue[T]{messages: NewQueue[T](), readers: NewQueue[chan T]()}
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
	i := mq.readers.push(ch)
	select {
	case <-ctx.Done():
		mq.readers.remove(i)
		close(ch)
		return nil, TimeoutError
	case v := <-ch:
		return &v, nil
	}
}
