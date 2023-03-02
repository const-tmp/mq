package queue

import (
	"testing"
)

const num = 1000000

func TestMessageQueueBefore(t *testing.T) {
	mq := New[int]()

	go func() {
		for i := 0; i < num; i++ {
			v, err := mq.Pop(0)
			if err != nil {
				t.Errorf("%d error: %s", i, err)
			}
			if v == nil {
				t.Errorf("%d result is nil", i)
			}
			if *v != i {
				t.Errorf("want %d got %d", i, *v)
			}
		}
	}()

	for i := 0; i < num; i++ {
		mq.Push(i)
	}
}

func TestMessageQueueAfter(t *testing.T) {
	mq := New[int]()

	for i := 0; i < num; i++ {
		mq.Push(i)
	}

	for i := 0; i < num; i++ {
		v, err := mq.Pop(0)
		if err != nil {
			t.Errorf("%d error: %s", i, err)
		}
		if v == nil {
			t.Errorf("%d result is nil", i)
		}
		if *v != i {
			t.Errorf("want %d got %d", i, *v)
		}
	}
}
