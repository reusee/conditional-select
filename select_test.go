package cselect

import "testing"

func TestSelect(t *testing.T) {
	selector := New()
	c1 := make(chan bool)
	c1pass := make(chan bool)
	c2 := make(chan int)
	c2pass := make(chan bool)
	c3 := make(chan string)
	c4 := make(chan int)
	c4pass := make(chan bool)
	c4ok := false

	selector.Add(c1, func(b bool) {
		close(c1pass)
	}, nil, nil)
	selector.Add(c2, func(i int) {
		close(c2pass)
	}, nil, nil)
	selector.Add(c3, nil, nil, func() string {
		return "foo"
	})
	selector.Add(c4, func() {
		close(c4pass)
	}, func() bool {
		return c4ok
	}, nil)

	go selector.Select()
	c1 <- true
	<-c1pass

	go selector.Select()
	c2 <- 42
	<-c2pass

	go selector.Select()
	if <-c3 != "foo" {
		t.Fatalf("c3")
	}

	go selector.Select()
	select {
	case c4 <- 1:
		t.Fatalf("c4 should not receive")
	default:
	}
	<-c3 // wake up selector

	c4ok = true
	go selector.Select()
	c4 <- 1
	<-c4pass
}
