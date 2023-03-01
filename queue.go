package queue

import (
	"context"
	"sync"
)

type Service struct {
	messages            []string
	appendMessage       chan string
	readers             []chan string
	appendReader        chan chan string
	messageMu, readerMu sync.RWMutex
}

func New(ctx context.Context) *Service {
	s := &Service{appendMessage: make(chan string), appendReader: make(chan chan string)}
	go s.appendReaderLoop(ctx)
	go s.appendMessageLoop(ctx)
	return s
}

func (s *Service) Put(v string) {
	if s.readersLen() > 0 {
		s.popReader() <- v
		return
	}
	s.appendMessage <- v
}

func (s *Service) Get() string {
	if s.messagesLen() > 0 {
		return s.popMessage()
	}
	ch := make(chan string)
	s.appendReader <- ch
	return <-ch
}

func (s *Service) appendMessageLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-s.appendMessage:
			s.messages = append(s.messages, v)
		}
	}
}

func (s *Service) appendReaderLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-s.appendReader:
			s.readers = append(s.readers, v)
		}
	}
}

func (s *Service) messagesLen() int {
	s.messageMu.RLock()
	defer s.messageMu.RUnlock()
	return len(s.messages)
}

func (s *Service) popMessage() string {
	s.messageMu.Lock()
	defer s.messageMu.Unlock()
	message := s.messages[0]
	s.messages = s.messages[1:]
	return message
}

func (s *Service) readersLen() int {
	s.readerMu.RLock()
	defer s.readerMu.RUnlock()
	return len(s.readers)
}

func (s *Service) popReader() chan string {
	s.readerMu.Lock()
	defer s.readerMu.Unlock()
	ch := s.readers[0]
	s.readers = s.readers[1:]
	return ch
}
