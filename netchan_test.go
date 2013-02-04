package netchan

import (
	"testing"
)

func Test1(t *testing.T) {

	done := make(chan bool)

	// server
	go func() {
		conn := Serve(":8080")
		ch := make(chan int)
		select {
		case err := <-conn.Connect(&ch, []byte("Test")):
			t.Errorf(err.Error())
		case n := <-ch:
			if n != 100 {
				t.Errorf("recieved value must be 100.")
			}

			done <- true
		}
	}()

	// send
	func() {
		conn, _ := Dial(":8080")
		ch := make(chan int)
		select {
		case err := <-conn.Connect(&ch, []byte("Test")):
			t.Errorf(err.Error())
		case ch <- 100:
		}
	}()

	<-done
}
