package queue

type (
	item[T any] struct {
		value   T
		expired bool
		next    *item[T]
	}

	fifo[T any] struct {
		first, last *item[T]
	}
)

func (l *fifo[T]) push(v T) *item[T] {
	i := &item[T]{value: v}
	if l.last == nil {
		l.first, l.last = i, i
		return i
	}
	l.last.next, l.last = i, i
	return i
}

func (l *fifo[T]) pop() (i *item[T]) {
	if l.first == nil {
		return
	}
	i, l.first = l.first, l.first.next
	for i.expired && l.first != nil {
		i, l.first = l.first, l.first.next
	}
	if i.expired {
		return nil
	}
	return
}
