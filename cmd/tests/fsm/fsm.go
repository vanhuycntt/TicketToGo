package main

import (
	"github.com/ryanfaerman/fsm"
)

type Thing struct {
	State fsm.State

	// our machine cache
	machine *fsm.Machine
}

// Add methods to comply with the fsm.Stater interface
func (t *Thing) CurrentState() fsm.State { return t.State }
func (t *Thing) SetState(s fsm.State)    { t.State = s }

// A helpful function that lets us apply arbitrary rulesets to this
// instances state machine without reallocating the machine. While not
// required, it's something I like to have.
func (t *Thing) Apply(r *fsm.Ruleset) *fsm.Machine {
	if t.machine == nil {
		t.machine = &fsm.Machine{Subject: t}
	}

	t.machine.Rules = r
	return t.machine
}

/*func main() {
	var err error

	some_thing := Thing{State: "pending"} // Our subject
	fmt.Println(some_thing)

	// Establish some rules for our FSM
	rules := fsm.Ruleset{}
	rules.AddTransition(fsm.T{"pending", "started"})
	rules.AddTransition(fsm.T{"started", "finished"})

	err = some_thing.Apply(&rules).Transition("started")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(some_thing)
}*/
