package queue

type (
	ListItem[T any] struct {
		value T
		next  *ListItem[T]
	}

	List[T any] struct {
		first, last *ListItem[T]
	}
)

func (l *List[T]) Push(v T) *ListItem[T] {
	i := &ListItem[T]{value: v}
	if l.last == nil {
		l.first, l.last = i, i
		return i
	}
	l.last.next, l.last = i, i
	return i
}

func (l *List[T]) Pop() (i *ListItem[T]) {
	if l.first == nil {
		return
	}
	i, l.first = l.first, l.first.next
	return
}

func (l *List[T]) Remove(i *ListItem[T]) {
	if l.first == i {
		l.first = nil
		return
	}
	current := l.first
	for current != nil {
		if current.next == i {
			current.next = i.next
			return
		}
		current = current.next
	}
}
