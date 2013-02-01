package netchan

import (
	"testing"
)

func Test1(t *testing.T) {

	// recieve
	go func() {
		ch := make(chan int)
		select {
		case err := <-Dial(&ch, nil, ":8080", ":9090", "Test"):
			if err != Done {
				t.Errorf(err.Error())
			}
		case n := <-ch:
			if n != 100 {
				t.Errorf("recieved value must be 100.")
			}
		}
	}()

	// send
	func() {
		ch := make(chan int)
		done := make(chan bool)
		select {
		case err := <-Dial(&ch, done, ":9090", ":8080", "Test"):
			if err != Done {
				t.Errorf(err.Error())
			}
		case ch <- 100:
			<-done
		}
	}()

}
