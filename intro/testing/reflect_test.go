package testing

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type Person interface {
	SayHello(name string)
	Run() string
}

type Hero struct {
	Name string
	Age int
	Speed int
}

func (h *Hero) SayHello(name string) {
	fmt.Printf("Hello %s, I am %s\n", name, h.Name)
}

func (h *Hero) Run() string {
	return fmt.Sprintf("I am running and my speed is %s\n", strconv.Itoa(h.Speed))
}

func TestReflect(t *testing.T) {
	typeOfHero := reflect.TypeOf(Hero{})
	fmt.Printf("Type of hero is %s, kind of hero is %s\n", typeOfHero, typeOfHero.Kind())
}

