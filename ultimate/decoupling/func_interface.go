package main

import "fmt"

type (
	Printer interface {
		Print()
	}
	Notifier interface {
		Notify()
	}
)
type (
	Duration int

	OutToConsole struct{}
)

func (o OutToConsole) Print() {
	fmt.Println(fmt.Sprintf("console output %p", &o))
}

func (d *Duration) Print() {
	fmt.Println(fmt.Sprintf("console output %d", *d))
}

func main() {
	o := OutToConsole{}

	outs := []Printer{o, &o}

	fmt.Println(fmt.Sprintf(" origin pointer %p", &o))
	fmt.Println(fmt.Sprintf(" 1st pointer %p", &outs[0]))
	fmt.Println(fmt.Sprintf(" 2st pointer %p", outs[1]))
	d := Duration(10)
	d.Print()
}
