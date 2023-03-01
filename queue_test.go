package queue

import (
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

const num = 100000

func TestServiceGetBefore(t *testing.T) {
	s := new(Service)

	go func() {
		for i := 0; i < num; i++ {
			if v := s.Get(); v != fmt.Sprint(i) {
				t.Errorf("want %d got %s", i, v)
			}
		}
	}()

	for i := 0; i < num; i++ {
		s.Put(fmt.Sprint(i))
	}
}

func TestServiceGetAfter(t *testing.T) {
	s := new(Service)

	for i := 0; i < num; i++ {
		s.Put(fmt.Sprint(i))
	}

	for i := 0; i < num; i++ {
		if v := s.Get(); v != fmt.Sprint(i) {
			t.Errorf("want %d got %s", i, v)
		}
	}
}
