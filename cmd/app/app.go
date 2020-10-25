package main

import "fmt"

func main() {
	/*	Filter([]Arr {&PersonArr{}, &PersonArr{}}, func(ar Arr)bool {
		return true
	})*/
	p := PersonArr{
		start: 10,
		end:   100,
	}

	arr := []PersonArr{p, p}
	p.start = 11
	p.end = 12

	for _, a := range arr {
		fmt.Println(a)
	}
}

type Predicater interface {
	Predicate(a Arr) bool
}

type Mapper interface {
	MapTo(a Arr) interface{}
}

func Filter(arr []Arr, predicater Predicater) []Arr {
	var rs []Arr
	for _, a := range arr {
		if !predicater.Predicate(a) {
			continue
		}
		if m, ok := predicater.(Mapper); ok {
			v := m.MapTo(a)
			rs = append(rs, v.(Arr))
		}
	}
	return rs
}

type Arr interface {
	Length() int
}

type PersonArr struct {
	start int
	end   int
}

func (p *PersonArr) Length() int {
	return 100
}
