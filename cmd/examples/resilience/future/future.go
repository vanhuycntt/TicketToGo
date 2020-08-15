package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Future interface {
	Result() (string, error)
}

type InnerFunction struct {
	doOnce sync.Once
	wg     sync.WaitGroup
	err    error
	rs     string

	errCh <-chan error
	rsCh  <-chan string
}

func (sf *InnerFunction) Result() (string, error) {
	sf.doOnce.Do(func() {
		sf.wg.Add(1)
		defer sf.wg.Done()
		sf.rs = <-sf.rsCh
		sf.err = <-sf.errCh

	})
	sf.wg.Wait()
	return sf.rs, sf.err
}

func SlowFunction(exCtx context.Context) Future {
	rsCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <-time.After(2 * time.Second):
			rsCh <- "by pass"
			errCh <- nil
		case <-exCtx.Done():
			rsCh <- ""
			errCh <- exCtx.Err()
		}

	}()

	return &InnerFunction{rsCh: rsCh, errCh: errCh}
}
func main() {
	exeCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ft := SlowFunction(exeCtx)
	msg, err := ft.Result()

	fmt.Println(msg)
	fmt.Println(err)

	msg, err = ft.Result()

	fmt.Println(msg)
	fmt.Println(err)
}
