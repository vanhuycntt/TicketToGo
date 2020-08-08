package main

import (
	"fmt"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(1)
}
func main() {

	/*n1 := Name{"abc", []byte("abc1")}
	n2 := Name{"abc", []byte("abc")}
	// using DeepEqual function to compare on incomparable types.
	isEqual := reflect.DeepEqual(n1, n2)
	//because slice []byte is incomparable, so the code line is compiled fail
	//isEqual = n1 == n2 //compile failed
	if isEqual {
		//fmt.Println("valid")
		fmt.Println("valid")
	}*/
	fmt.Println("start")

	go saySomething()
	saySomething()
	fmt.Println("end")
}
func saySomething() {

	for i := 0; i < 10; i++ {
		runtime.Gosched()
		fmt.Println(" hello ", i)
	}

}

type Name struct {
	name string
	raws []byte
}
