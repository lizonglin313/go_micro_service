package test

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	c := make(chan string, 2)
	c <- "hello"
	c <- "world"

	time.AfterFunc(time.Second*3, func() {
		close(c)
	})

	for e := range c {
		fmt.Printf("e is : %s\n", e)
	}
}
