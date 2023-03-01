package queue

import (
	"sync"
)

type Service struct {
	messages            []string
	readers             []chan string
	messageMu, readerMu sync.RWMutex
}

func (s *Service) Put(message string) {
	if s.readersLen() > 0 {
		s.popReader() <- message
		return
	}
	s.appendMessage(message)
}

func (s *Service) Get() string {
	if s.messagesLen() > 0 {
		return s.popMessage()
	}
	ch := make(chan string)
	s.appendReader(ch)
	return <-ch
}

func (s *Service) messagesLen() int {
	s.messageMu.RLock()
	defer s.messageMu.RUnlock()
	return len(s.messages)
}

func (s *Service) appendMessage(message string) {
	s.messageMu.Lock()
	defer s.messageMu.Unlock()
	s.messages = append(s.messages, message)
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

func (s *Service) appendReader(ch chan string) {
	s.readerMu.Lock()
	defer s.readerMu.Unlock()
	s.readers = append(s.readers, ch)
}

func (s *Service) popReader() chan string {
	s.readerMu.Lock()
	defer s.readerMu.Unlock()
	ch := s.readers[0]
	s.readers = s.readers[1:]
	return ch
}
