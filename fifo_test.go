package mq

import "testing"

func TestList(t *testing.T) {
	l := new(fifo[int])
	for i := 0; i < num; i++ {
		l.push(i)
	}
	for i := 0; i < num; i++ {
		if item := l.pop(); item == nil || item.value != i {
			t.Error()
		}
	}
}

func TestFifoExpiration(t *testing.T) {
	f := new(fifo[int])
	for i := 0; i < num; i++ {
		item := f.push(i)
		if i%2 != 0 {
			item.expired = true
		}
	}
	for i := 0; i < num; i++ {
		if item := f.pop(); item != nil && item.expired {
			t.Errorf("item %d is expired", i)
		}
	}
}
