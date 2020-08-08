package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	syncCall := SysCall{
		mutex: &sync.Mutex{},
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println(syncCall.Add(10))
	}()
	wg.Wait()
}

type SysCall struct {
	mutex   *sync.Mutex
	counter int
}

func (s *SysCall) Add(incr int) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.checkPerms() {
		s.counter = s.counter + incr
	}
	return s.counter
}

func (s *SysCall) checkPerms() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.counter > 0 {
		return true
	}
	return false
}
