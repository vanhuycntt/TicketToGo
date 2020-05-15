package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	rateLimiter := rate.NewLimiter(1, 3)

	allow := rateLimiter.Wait(context.Background())
	fmt.Println(allow)
	time.Sleep(time.Second)
	allow = rateLimiter.Wait(context.Background())

	fmt.Println(allow)
	time.Sleep(time.Second)
	allow = rateLimiter.Wait(context.Background())
	fmt.Println(allow)
	time.Sleep(time.Second)
	allow = rateLimiter.Wait(context.Background())
	fmt.Println(allow)
	time.Sleep(time.Second)
}
