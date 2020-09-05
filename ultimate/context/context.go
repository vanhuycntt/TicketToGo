package main

import (
	"context"
	"fmt"
	"time"
)

type ID struct {
	id   int32
	name string
}
type Pattern struct {
	header string
	body   string
	footer string
}

var stores = map[ID]Pattern{
	ID{10, "abc"}: {header: "Dunning money", body: "Let me know", footer: "Get it"},
	ID{11, "abc"}: {header: "Dunning money", body: "Let me know", footer: "Get it 01"},
}

func Get(ctx context.Context, id ID) <-chan Pattern {
	pCh := make(chan Pattern)
	go func() {
		time.Sleep(2 * time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("received timeout")
			close(pCh)
		default:
			pCh <- stores[id]
		}
	}()
	return pCh
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chArr := []<-chan Pattern{Get(ctx, ID{12, "abc"}), Get(ctx, ID{12, "abc"})}
	for _, ch := range chArr {
		p, ok := <-ch
		if ok {
			empty := Pattern{}
			if p != empty {
				fmt.Println(p)
			} else {
				fmt.Println("no value")
			}
		} else {
			fmt.Println("timout")
		}
	}

}
