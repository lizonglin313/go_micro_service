package testing

import (
	"fmt"
	"testing"
)

type Printer interface {
	Print(interface{})
}

type FuncCaller func(p interface{})

func (f FuncCaller) Print(i interface{}) {
	f(i)
}

func TestInterface(t *testing.T) {
	var printer Printer
	printer = FuncCaller(func(p interface{}) {
		fmt.Println(p)
	})
	printer.Print("good")
}