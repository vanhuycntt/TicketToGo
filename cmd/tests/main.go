package main

import "fmt"

func main() {
	temp := struct {
		F uint32
		C uint32
	}{
		15,
		30,
	}
	fmt.Println(fmt.Sprintf("%p", &temp))
	fmt.Println(fmt.Sprintf("%p", &(temp).F))

	temp01 := temp
	fmt.Println(fmt.Sprintf("temp %p", &temp01))
	fmt.Println(fmt.Sprintf("temp %p", &(temp01).F))

}
