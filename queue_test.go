package mq

import (
	"testing"
)

const num = 10000

func TestMessageQueue_Push(t *testing.T) {
	mq := new(MQ[int])
	for i := 0; i < num; i++ {
		mq.Push(i)
	}
}

func TestMessageQueueBefore(t *testing.T) {
	var (
		v   int
		err error
		mq  = new(MQ[int])
	)

	go func() {
		for i := 1; i < num; i++ {
			v, err = mq.Pop(0)
			if err != nil {
				t.Errorf("%d error: %s", i, err)
			}
			if v == 0 {
				t.Errorf("%d result is nil", i)
			}
			if v != i {
				t.Errorf("want %d got %d", i, v)
			}
		}
	}()

	for i := 1; i < num; i++ {
		mq.Push(i)
	}
}

func TestMessageQueueAfter(t *testing.T) {
	mq := new(MQ[int])

	for i := 1; i < num; i++ {
		mq.Push(i)
	}

	for i := 1; i < num; i++ {
		v, err := mq.Pop(0)
		if err != nil {
			t.Errorf("%d error: %s", i, err)
		}
		if v == 0 {
			t.Errorf("%d result is nil", i)
		}
		if v != i {
			t.Errorf("want %d got %d", i, v)
		}
	}
}

func BenchmarkMessageQueue_Push(b *testing.B) {
	mq := new(MQ[int])
	for i := 0; i < b.N; i++ {
		mq.Push(i)
	}
}

func BenchmarkMessageQueue_Pop(b *testing.B) {
	mq := new(MQ[int])
	for i := 0; i < b.N; i++ {
		mq.Push(i)
	}
	for i := 0; i < b.N; i++ {
		mq.Pop(0)
	}
}

func BenchmarkMessageQueue_Pop2(b *testing.B) {
	mq := new(MQ[int])
	for i := 0; i < b.N; i++ {
		mq.Pop(1)
	}
}
