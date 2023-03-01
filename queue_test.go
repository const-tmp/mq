package queue

import (
	"context"
	"fmt"
	"testing"
)

func TestChan(t *testing.T) {
	c := make(chan int)
	go recv(c)
	send(c, 1)
}

func recv(c chan int) {
	fmt.Println("recv before")
	fmt.Println("recv:", <-c)
	fmt.Println("recv after")
}

func send(c chan int, v int) {
	fmt.Println("send before")
	c <- v
	fmt.Println("send after")
}

func TestService(t *testing.T) {
	num := 1000
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := New(ctx)

	for i := 0; i < num; i++ {
		go func(n int) {
			if v := s.Get(); v != fmt.Sprint(n) {
				t.Errorf("want %d got %s", n, v)
			}
		}(i)
		//time.Sleep(time.Millisecond)
	}

	for i := 0; i < num; i++ {
		s.Put(fmt.Sprint(i))
	}
}
