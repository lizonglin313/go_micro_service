package main

import (
	"fmt"
)

func SliceRise(s []int) {
	s = append(s, 0)
	for i:= range s {
		s[i]++
	}
	fmt.Println(s)
}

func main() {


	orderLen := 5
	order := make([]uint16, orderLen*2)
	order = []uint16{1,2,3,4,5,6,7,8,9,10}


	pollorder := order[:orderLen:orderLen]
	lockorder := order[orderLen:][:orderLen:orderLen]

	fmt.Println(pollorder)
	fmt.Println(lockorder)

}

func fibo() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}