package mq

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

func (f *fifo[T]) push(v T) *item[T] {
	i := &item[T]{value: v}
	if f.last == nil {
		f.first, f.last = i, i
		return i
	}
	f.last.next, f.last = i, i
	return i
}

func (f *fifo[T]) pop() (i *item[T]) {
	if f.first == nil {
		return
	}
	i, f.first = f.first, f.first.next
	for i.expired && f.first != nil {
		i, f.first = f.first, f.first.next
	}
	if i.expired {
		return nil
	}
	return
}

func (f *fifo[T]) pop2() (i *item[T]) {
	for f.first != nil {
		if !f.first.expired {
			i, f.first = f.first, f.first.next
			break
		}
		f.first = f.first.next
	}
	return
}

func (f *fifo[T]) pop3() (i *item[T]) {
	for f.first != nil && f.first.expired {
		f.first = f.first.next
	}
	if f.first == nil {
		return
	}
	i, f.first = f.first, f.first.next
	return
}
