package channel_shared_seats

import (
	"TicketToGo"
	"context"
	"errors"
)

type RequestType string

const (
	HoldType   RequestType = "HOLD"
	UnHoldType RequestType = "UNHOLD"
)

type (
	Request struct {
		typeReq    RequestType
		seatNum    int
		acceptChan chan interface{}
	}

	SeatChartCoordinator struct {
		seats  []*TicketsToGo.Seat
		length int

		head int
		tail int

		requestPool chan *Request
	}
)

func newSeatChartCoordinator(numSeats int) *SeatChartCoordinator {
	seatChart := &SeatChartCoordinator{
		seats:  make([]*TicketsToGo.Seat, numSeats),
		length: numSeats,
		head:   0,
		tail:   numSeats - 1,

		requestPool: make(chan *Request, 10),
	}
	go seatChart.coordinate(context.Background())

	return seatChart
}

func (sc *SeatChartCoordinator) Hold(num int) chan interface{} {
	holdSeats := make(chan interface{}, 1)
	req := Request{
		typeReq:    HoldType,
		seatNum:    num,
		acceptChan: holdSeats,
	}
	sc.requestPool <- &req

	return holdSeats
}

func (sc *SeatChartCoordinator) UnHold(restoredSeats []int) {
	reqChan := make(chan interface{}, 1)
	reqChan <- restoredSeats
	req := Request{
		typeReq:    UnHoldType,
		acceptChan: reqChan,
	}
	sc.requestPool <- &req
}

func (sc *SeatChartCoordinator) coordinate(ctx context.Context) {
escapeLoop:
	for {
		select {
		case req := <-sc.requestPool:
			if req.typeReq == HoldType {
				sc.hold(req)
			}
			if req.typeReq == UnHoldType {
				sc.unhold(req)
			}
		case <-ctx.Done():
			break escapeLoop
		default:
		}
	}
}

func (sc *SeatChartCoordinator) hold(req *Request) {
	if sc.head == sc.tail || (sc.head+req.seatNum) > sc.tail {
		req.acceptChan <- errors.New("Not Enough Seats")
		return
	}

	from, to := sc.head, sc.head+req.seatNum
	for s := from; s < to; s++ {
		seatObj := sc.seats[s]

		if seatObj == nil {
			pos := s % sc.length
			sc.seats[s] = &TicketsToGo.Seat{
				PlacedFlag: TicketsToGo.HOLD,
				Position:   pos,
			}
			continue
		}

		if seatObj.PlacedFlag == TicketsToGo.UNHOLD {
			seatObj.PlacedFlag = TicketsToGo.HOLD
		}
	}
	var posSeats []int
	for _, s := range sc.seats[from:to] {
		posSeats = append(posSeats, s.Position)
	}
	sc.head += req.seatNum
	req.acceptChan <- posSeats
}

func (sc *SeatChartCoordinator) unhold(req *Request) {
	var compensateSeats []*TicketsToGo.Seat
escapeLoop:
	for {
		select {
		case seatsArr := <-req.acceptChan:
			positions := seatsArr.([]int)
			compensateSeats = make([]*TicketsToGo.Seat, len(positions))
			for idx, pos := range positions {
				compensateSeats[idx] = &TicketsToGo.Seat{
					PlacedFlag: TicketsToGo.UNHOLD,
					Position:   pos,
				}
			}
			break escapeLoop
		}
	}

	sc.seats = append(sc.seats, compensateSeats...)

	sc.tail += len(compensateSeats)
}
func (sc *SeatChartCoordinator) Collect() []*TicketsToGo.Seat {
	rsSeats := make(map[int]*TicketsToGo.Seat)
	for idx, s := range sc.seats {
		if s == nil {
			rsSeats[idx] = &TicketsToGo.Seat{
				PlacedFlag: TicketsToGo.UNHOLD,
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
	var seats []*TicketsToGo.Seat
	for _, val := range rsSeats {
		seats = append(seats, val)
	}
	return seats
}
