package channel_shared_seats

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestMSenderNReceiver(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	nRec := 100
	mSend := 10

	syncWait := sync.WaitGroup{}
	syncWait.Add(nRec)

	dataCh := make(chan int)
	stopCh := make(chan struct{})
	toStop := make(chan string, 1)

	go func() {
		str := <-toStop
		fmt.Println("stop message: " + str)
		close(stopCh)
	}()

	//senders
	for i := 0; i < mSend; i++ {
		go func(id int) {
			for {
				val := rand.Int()
				deadVal := rand.Intn(10000)
				if deadVal == 9999 {

					select {
					case toStop <- fmt.Sprintf("sender#%d send %d", id, deadVal):
					default:
					}
					return

				}

				select {
				case <-stopCh:
					return
				default:

				}
				select {
				case <-stopCh:
					return
				case dataCh <- val:

				}
			}

		}(i)
	}
	//receivers
	for i := 0; i < nRec; i++ {
		go func(id int) {
			defer syncWait.Done()
			for {
				select {
				case <-stopCh:
					return
				default:

				}
				select {
				case <-stopCh:
					return
				case msg := <-dataCh:
					fmt.Println(fmt.Sprintf("receiver#%d: %d", id, msg))
				}
			}

		}(i)
	}
	syncWait.Wait()
}
