package client

import (
	"testing"
	"gopkg.in/go-playground/pool.v3"
	"fmt"
	"time"
)

func TestSendEmail(t *testing.T) {
	p := pool.NewLimited(1)
	defer p.Close()

	batch := p.Batch()

	// for max speed Queue in another goroutine
	// but it is not required, just can't start reading results
	// until all items are Queued.

	go func() {
		for i := 0; i < 10; i++ {
			batch.Queue(SendEmail("email content"))
		}

		// DO NOT FORGET THIS OR GOROUTINES WILL DEADLOCK
		// if calling Cancel() it calles QueueComplete() internally
		batch.QueueComplete()
	}()
	time.Sleep(10 * time.Second)
	fmt.Println("111")
	for email := range batch.Results() {

		if err := email.Error(); err != nil {
			// handle error
			// maybe call batch.Cancel()
		}

		// use return value
		fmt.Println(email.Value().(bool))
	}
}
