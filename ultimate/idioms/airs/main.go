package main

import "fmt"

func main() {
	var mv Mover
	mv = Humman{}

	if st, ok := mv.(Steper); ok {
		fmt.Println("step by 10 m")
		st.StepBy(10)
	}
	var nm NextMover
	nm = Humman{}
	accept(nm)

}

func accept(mv Mover) {
	fmt.Println(mv)
}

var _ Mover = (*Humman)(nil)

type (
	Mover interface {
		Move(length int)
	}

	Steper interface {
		StepBy(number int)
	}

	Runner interface {
		Run()
	}

	MoverRuner interface {
		Mover
		Runner
	}

	NextMover interface {
		Move(length int)
	}
)

type (
	Moving struct {
	}
	Steping struct {
	}

	Humman struct {
		Moving
		Steping
	}

	Animal struct {
	}
)

func (mv Moving) Move(length int) {

}

func (st Steping) StepBy(number int) {

}
