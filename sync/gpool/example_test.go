package gpool

import (
	"fmt"
)

func ExampleNew() {
	//we create a poll named "p" with 3 workers
	p := New(3)
	//after New, you can submit a function as a task, you can repeat Submit() many times anywhere as you need.
	a := make(chan bool)
	p.Submit(func() {
		a <- true
	})
	fmt.Println(<-a)
}
