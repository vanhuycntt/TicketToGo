package main

import (
	"flag"
	"fmt"
)

func main() {
	var temp = Temp{
		f: 0.1,
	}
	flag.Var(&temp, "temp", "your temperature")

	flag.Parse()

	fmt.Println(temp.String())
}

type Temp struct {
	f float64
	d string
}

func (t *Temp) Set(str string) error {
	_, err := fmt.Sscanf(str, "%f%s", &t.f, &t.d)
	return err
}

func (t *Temp) String() string {
	return fmt.Sprintf("%f%s", t.f, t.d)
}
