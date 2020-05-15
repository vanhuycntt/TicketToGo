package main

import (
	"fmt"
	"math"
	"time"
)

const (
	expirationPeriod = 5
)

func main() {
	//videoUUID := uuid.New().ID()
	//fmt.Println("uuid", videoUUID)
	//cacheKey := time.Now().Unix() + expirationPeriod
	//fmt.Println(cacheKey)
	//jitter := videoUUID % uint32(expirationPeriod)
	//fmt.Println("jitter: ", jitter)
	//cacheWindowSize := (time.Now().Unix() + int64(jitter))  / expirationPeriod
	//fmt.Println(cacheWindowSize)
	//
	//time.Sleep(6 * time.Second)
	//cacheWindowSize = (time.Now().Unix() + int64(jitter)) / expirationPeriod
	//fmt.Println(cacheWindowSize)
	//uuidStr := uuid.New().String()
	//fmt.Println(uuidStr)
	//encodeStr := base32.StdEncoding.EncodeToString([]byte(uuidStr))
	//fmt.Println(encodeStr)
	//byteDecoded, _ := base32.StdEncoding.DecodeString(encodeStr)
	//fmt.Println(string(byteDecoded))

	/*var sync sync.WaitGroup
	values := []int {1,2,3}
	errCh := make(chan error)
	doneCh := make(chan struct{})

	for _, v := range values {
		sync.Add(1)
		go func(idx int) {
			defer sync.Done()
			errCh <- fmt.Errorf("error on %v", idx)
		}(v)
	}
	go func() {
		sync.Wait()
		close(doneCh)
	}()
	select {
	case err := <- errCh:
		fmt.Println(err)
	case <- doneCh:
		fmt.Println("done")
	}*/
	//fmt.Println(Eod(time.Now()))

	currentTimeStamp := time.Now().Unix()
	ts := time.Now().Add(10 * time.Second).Unix()
	ttl := 5 * time.Second
	fmt.Println(math.Abs(float64(currentTimeStamp - ts)))
	fmt.Println(math.Abs(ttl.Seconds()))

	if math.Abs(float64(currentTimeStamp-ts)) > ttl.Seconds() {
		fmt.Println("over time to live")
	}

}
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
func Eod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 11, 59, 59, 0, t.Location())
}
