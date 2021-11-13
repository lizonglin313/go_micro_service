package testing

import (
	"fmt"
	"testing"
)

func createCounter(initial int) func() int {
	if initial <0 {
		initial = 0
	}

	return func() int {
		initial++
		return initial
	}
}

func TestClosure(t *testing.T) {
	c1 := createCounter(1)
	fmt.Println(c1())
	fmt.Println(c1())

	c2 := createCounter(100)
	fmt.Println(c2())
	fmt.Println(c1())
}
