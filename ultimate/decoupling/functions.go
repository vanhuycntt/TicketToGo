package main

type Attention struct {
	ValidFunc func() bool
}

type AttendTo func(number int32) int

func AttendMembers() {
	var members = AttendTo(func(number int32) int {
		return int(number) - 10
	})
	members(100)
}
