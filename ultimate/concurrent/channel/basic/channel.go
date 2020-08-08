package basic

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var ch = make(chan struct{})
	//var tempCh = make(chan struct{})

	syncW := sync.WaitGroup{}
	syncW.Add(2)
	/*go func() {
		defer syncW.Done()
		fmt.Println("g 01")
		_, ok :=<- ch
		if !ok {
			fmt.Println("g 01=> closed")
		}
	}()

	go func() {
		defer syncW.Done()
		fmt.Println("g 02")
		_, ok := <- ch
		if !ok {
			fmt.Println("g 02=> closed")
		}
		for _ = range ch {
			fmt.Println("g 02=> closed")
		}
	}()
	*/

	go func() {
		defer syncW.Done()
		time.Sleep(time.Second)
		/*select {
			case <- ch:
				fmt.Println("ch")
			case <- tempCh:
				fmt.Println("tempCh")
		}*/
		for v := range ch {
			fmt.Println("ch01", v)
		}
	}()

	go func() {
		defer syncW.Done()
		time.Sleep(time.Second)
		/*select {
			case <- ch:
				fmt.Println("ch")
			case <- tempCh:
				fmt.Println("tempCh")
		}*/
		for v := range ch {
			fmt.Println("ch02", v)
		}
	}()

	ch <- struct{}{}
	close(ch)
	syncW.Wait()
	fmt.Println("end")

}
