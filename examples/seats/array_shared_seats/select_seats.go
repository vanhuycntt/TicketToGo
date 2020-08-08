package array_shared_seats

import (
	"github.com/pkg/errors"
	"sync"
)

type HoldFlag int

const (
	HOLD HoldFlag = iota
	UNHOLD
)

type SeatsChart struct {
	seats  []*Seat
	mux    *sync.RWMutex
	length int

	head int
	tail int
}

type Seat struct {
	PlacedFlag HoldFlag
	Position   int
}

func NewSeatsChart(num int) *SeatsChart {
	return &SeatsChart{
		length: num,
		seats:  make([]*Seat, num),
		mux:    &sync.RWMutex{},

		head: 0,
		tail: num - 1,
	}
}

func (sc *SeatsChart) Hold(placedSeats int) ([]int, error) {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	if sc.head == sc.tail || (sc.head+placedSeats) > sc.tail {
		return nil, errors.New("Not Enough Seats")
	}

	from, to := sc.head, sc.head+placedSeats

	for s := from; s < to; s++ {
		seatObj := sc.seats[s]

		if seatObj == nil {
			pos := s % sc.length
			sc.seats[s] = &Seat{
				PlacedFlag: HOLD,
				Position:   pos,
			}
			continue
		}

		if seatObj.PlacedFlag == UNHOLD {
			seatObj.PlacedFlag = HOLD
		}

	}
	var posSeats []int
	for _, s := range sc.seats[from:to] {
		posSeats = append(posSeats, s.Position)
	}
	sc.head += placedSeats
	return posSeats, nil
}

func (sc *SeatsChart) UnHold(unplacedSeats []int) error {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	var compensateSeats []*Seat
	for _, pos := range unplacedSeats {
		compensateSeats = append(compensateSeats, &Seat{
			PlacedFlag: UNHOLD,
			Position:   pos,
		})
	}
	sc.seats = append(sc.seats, compensateSeats...)

	sc.tail += len(unplacedSeats)
	return nil
}

func (sc *SeatsChart) Collect() []*Seat {
	sc.mux.RLock()
	defer sc.mux.RUnlock()
	rsSeats := make(map[int]*Seat)
	for idx, s := range sc.seats {
		if s == nil {
			rsSeats[idx] = &Seat{
				PlacedFlag: UNHOLD,
				Position:   idx,
			}
			continue
		}
		if val, ok := rsSeats[s.Position]; ok {
			val.PlacedFlag = s.PlacedFlag
			continue
		}
		rsSeats[s.Position] = s
	}
	var seats []*Seat
	for _, val := range rsSeats {
		seats = append(seats, val)
	}
	return seats
}
