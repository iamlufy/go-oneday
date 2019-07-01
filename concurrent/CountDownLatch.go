package concurrent

import (
	"fmt"
	"time"
)

type CountDownLatch struct {
	c chan interface{}
}

func NewCountDownLatch(count int) *CountDownLatch {
	return &CountDownLatch{
		c: make(chan interface{}, count),
	}
}
func (c *CountDownLatch) countDown() {
	c.c <- true
}
func (c *CountDownLatch) await() bool {
	for len(c.c) < cap(c.c) {
		fmt.Println("wait")
	}
	return true
}
func (c *CountDownLatch) awaitWithTime(d time.Duration) {
	start := time.Now()
	for len(c.c) < cap(c.c) {
		fmt.Println("wait")
		now := time.Now()
		if start.Add(d).Before(now)   {
			fmt.Println("time out")
			break
		}
	}
}

func main() {
	c := NewCountDownLatch(5)
	go func() { c.countDown() }()
	go func() { c.countDown() }()
	go func() { c.countDown() }()
	go func() { c.countDown() }()
	go func() {
		time.Sleep(3e6)
		c.countDown()
	}()
	c.awaitWithTime(4e6)
}
