package main

import (
	"encoding/json"
	"fmt"
	"github.com/looplab/fsm"
	"runtime"
)

type Door struct {
	To  string
	FSM *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.enterState(e) },
		},
	)

	return d
}

func (d *Door) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}

type VhEventMapper struct {
	EventID              uint            `json:"event_id"`
	TicketTypeID         uint            `json:"ticket_type_id"`
	ReceivedTicketTypeID uint            `json:"received_ticket_type_id"`
	ShippingInfo         json.RawMessage `json:"shipping_info"`
}

type ShippingInfo struct {
	ReceiverAddr  string `json:"receiver_addr"`
	ReceiverPhone string `json:"receiver_phone"`
}

func main() {
	/*door := NewDoor("heaven")

	err := door.FSM.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	err = door.FSM.Event("close")
	if err != nil {
		fmt.Println(err)
	}*/

	/*sysEndChan := make(chan os.Signal, 1)
	signal.Notify(sysEndChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan interface{})
	go func() {
		select {
		case sig := <-sysEndChan:
			fmt.Println(fmt.Sprintf("end by signal: %v", sig))
		case msg := <-done:
			fmt.Println(fmt.Sprintf("end by manual stop %v", msg ))
		}
		fmt.Println("close all resources")
	}()

	timer :=time.AfterFunc(10 * time.Second, func() {
		close(done)

	})
	time.Sleep(5 * time.Second)
	timer.Stop()*/
	/*jsonStr := `{
					"event_id":41752914,
					"ticket_type_id":91091319,
					"received_ticket_type_id":6925,
					"shipping_info" : "15 Cao Hung, Q3, tpHCM"
				}`

	var eventMapper VhEventMapper
	_ = json.Unmarshal([]byte(jsonStr), &eventMapper)
	var shippingInfo ShippingInfo
	err := json.Unmarshal(eventMapper.ShippingInfo, &shippingInfo)
	if err != nil  {
		unmarshalErr , ok := err.(*json.UnmarshalTypeError)
		if ok {
			fmt.Println(unmarshalErr)
		}
	}
	fmt.Println(shippingInfo.ReceiverAddr)
	fmt.Println(shippingInfo.ReceiverPhone)*/

	fmt.Println(runtime.GOMAXPROCS(-1))
	runtime.GOMAXPROCS(100)
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(-1))

}
