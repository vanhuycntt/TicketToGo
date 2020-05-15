package channel_shared_seats

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/rand"
	"time"

	"sync"
	"testing"
)

func RunMultiUserSelectSeatChartCoordinator(numUser int, numSeats int) {
	seatsChart := newSeatChartCoordinator(numSeats)

	var sWait sync.WaitGroup

	for i := 0; i < numUser; i++ {
		sWait.Add(1)
		go func(id int) {
			defer sWait.Done()
			num := rand.IntnRange(1, 5)
			resChan := seatsChart.Hold(num)
			r := <-resChan
			if _, ok := r.(error); ok {
				fmt.Println(fmt.Sprintf("%d hold %d error: %v", id, num, r.(error)))
				return
			}
			posSeats := r.([]int)
			fmt.Println(fmt.Sprintf("%d hold seats: %v", id, posSeats))
			if id%2 == 0 {
				seatsChart.UnHold(posSeats)

				fmt.Println(fmt.Sprintf("%d unhold seats: %v", id, posSeats))
			}

		}(i)
	}
	time.Sleep(time.Second * 2)
	sWait.Wait()
	for _, s := range seatsChart.Collect() {
		fmt.Println(fmt.Sprintf("SeatChart result %v", s))
	}
}

func TestSeatsChartCoordinator(b *testing.T) {
	for i := 0; i < 10; i++ {
		RunMultiUserSelectSeatChartCoordinator(100, 100)
	}
}

type ItemData interface{}

type RegularItem string

type DeadItem string

func TestClosedChannel(b *testing.T) {
	reqChan := make(chan ItemData, 2)
	syncWait := sync.WaitGroup{}
	syncWait.Add(1)

	go func() {
		defer syncWait.Done()
	escapeInfinityLoop:
		for {
			select {
			case r := <-reqChan:
				switch r.(type) {
				case RegularItem:
					fmt.Println(fmt.Sprintf("%v", r))
				case DeadItem:
					break escapeInfinityLoop
				}
			default:
			}
		}
	}()

	reqChan <- RegularItem("Regular 01")
	reqChan <- RegularItem("Regular 02")
	reqChan <- DeadItem("Dead")

	syncWait.Wait()

}
