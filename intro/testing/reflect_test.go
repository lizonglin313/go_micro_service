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
	Name  string
	Age   int
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

	for i := 0; i < typeOfHero.NumField(); i++ {
		fmt.Printf("field name is %s, type is %s, kind is %s\n",
			typeOfHero.Field(i).Name,
			typeOfHero.Field(i).Type,
			typeOfHero.Field(i).Type.Kind())
	}

	nameField, _ := typeOfHero.FieldByName("Name")
	fmt.Printf("method name is %s, type is %s, kind is %s\n",
		nameField.Name,
		nameField.Type,
		nameField.Type.Kind())
}

func TestGetFunctionUsingReflect(t *testing.T) {
	var person Person = &Hero{}
	typeofPerson  := reflect.TypeOf(person)

	for i := 0; i < typeofPerson.NumMethod(); i++ {
		fmt.Printf("method name is %s, type is %s, kind is %s\n",
			typeofPerson.Method(i).Name,
			typeofPerson.Method(i).Type,
			typeofPerson.Method(i).Type.Kind())
	}

	runMethod, _ := typeofPerson.MethodByName("Run")
	fmt.Printf("field name is %s, type is %s, kind is %s\n",
		runMethod.Name,
		runMethod.Type,
		runMethod.Type.Kind())
}

// .Elem() 用来解引用
func TestReflectValue(t *testing.T) {
	name := "小明"
	valueofName := reflect.ValueOf(&name)
	valueofName.Elem().Set(reflect.ValueOf("小红"))
	fmt.Println(name)
}

func TestCallFunc(t *testing.T) {
	var person Person = &Hero{
		Name: "小红",
		Speed: 12,
	}

	valueofPerson := reflect.ValueOf(person)
	sayHello := valueofPerson.MethodByName("SayHello")
	sayHello.Call([]reflect.Value{reflect.ValueOf("小张")})

	run := valueofPerson.MethodByName("Run")
	result := run.Call([]reflect.Value{})
	fmt.Printf("Result of run is %s", result[0])	// 获取函数结果

	// 如果用 TypeOf 会丢失接收器，需要显式传入
	valPerson := reflect.TypeOf(person)
	say, _ := valPerson.MethodByName("SayHello")
	say.Func.Call([]reflect.Value{reflect.ValueOf(person), reflect.ValueOf("dkcbbb")})
}