package testing

import (
	"fmt"
	"testing"
)

func pro(input string, processer func(str string)) {
	processer(input)
}

func TestCallBack(t *testing.T) {
	processer := func(str string) {
		for _, v := range str {
			fmt.Printf("%c ", v)
		}
	}
	pro("dkcbbb", processer)
}
