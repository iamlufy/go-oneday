package concurrent

import (
	"fmt"
	"time"
	"go-oneday/atomic"
)

type CyclicBarrier struct {
	ch     chan interface{}
	number atomic.Volatile
}

func NewCyclicBarrier(number uint) *CyclicBarrier {
	c := &CyclicBarrier{
		ch:     make(chan interface{}, number),
		number: atomic.NewVolatile(number),
	}
	//c.ch <- number
	c.number = atomic.NewVolatile(number)
	return c
}
func (c *CyclicBarrier) await() {
	c.ch <- true
	for len(c.ch) < cap(c.ch) {

	}
}
func main() {
	c := NewCyclicBarrier(4)
	go func() { fmt.Println("executing mission 1"); c.await(); fmt.Println("continue executing mission 1") }()
	go func() { fmt.Println("executing mission 2"); c.await(); fmt.Println("continue executing mission 2") }()
	go func() { fmt.Println("executing mission 3"); c.await(); fmt.Println("continue executing mission 3") }()
	go func() { fmt.Println("executing mission 4"); c.await(); fmt.Println("continue executing mission 4") }()
	time.Sleep(300e9)
}
