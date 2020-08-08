package array_shared_seats

import (
	"TicketToGo"
	"fmt"
	"k8s.io/apimachinery/pkg/util/rand"
	"runtime"

	"sync"
	"testing"
)

func RunMultiUserSelectSeatChart(numUser int, numSeats int) {
	seatsChart := TicketsToGo.NewSeatsChart(numSeats)

	var sWait sync.WaitGroup

	for i := 0; i < numUser; i++ {
		sWait.Add(1)
		go func(id int) {
			defer sWait.Done()
			num := rand.IntnRange(1, 5)
			posSeats, err := seatsChart.Hold(num)
			if err != nil {
				fmt.Println(fmt.Sprintf("%d hold error: %v", id, err))
				return
			}
			fmt.Println(fmt.Sprintf("%d hold seats: %v", id, posSeats))
			if id%2 == 0 {
				err = seatsChart.UnHold(posSeats)
				if err != nil {
					fmt.Println(fmt.Sprintf("%d unhold error: %v", id, err))
					return
				}
				fmt.Println(fmt.Sprintf("%d unhold seats: %v", id, posSeats))
			}

		}(i)
	}

	sWait.Wait()
	for _, s := range seatsChart.Collect() {
		fmt.Println(fmt.Sprintf("SeatChart result %v", s))
	}
}
func BenchmarkSeatsChart(b *testing.B) {
	n := runtime.GOMAXPROCS(0)

	for i := 0; i < b.N; i++ {
		RunMultiUserSelectSeatChart(n, 100)
	}
}
