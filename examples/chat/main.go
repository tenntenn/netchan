package main

import (
	"flag"
	"fmt"
	//"github.com/tenntenn/netchan"
	"../../../netchan"
	"bufio"
	"os"
	"time"
)

// program arguments.
var (
	ownHost    = flag.String("h", ":8080", "Own host address")
	remoteHost = flag.String("r", ":9090", "Remote host address")
)

func init() {
	flag.Parse()
}

func run(closed chan bool) {

	ch := make(chan string)
	done := make(chan bool)
	errCh := netchan.Dial(&ch, done, *ownHost, *remoteHost, "chat")
	r := bufio.NewReader(os.Stdin)
	for {
		select {
		case err := <-errCh:
			if err != nil {
				closed <- true
				return
			}
		case msg := <-ch:
			fmt.Println("Remote host says: ", msg)
		case <-closed:
			return
		default:

			fmt.Print(">")
			msg, _, err := r.ReadLine()
			if err != nil {
				fmt.Print(">")
				continue
			}
			ch <- string(msg)
			select {
			case <-done:
				fmt.Println("You said: ", msg)
			case <-time.After(time.Duration(5) * time.Second):
				fmt.Println("time out")
				closed <- true
				break
			}
		}
	}
}

func main() {
	closed := make(chan bool)
	go run(closed)
	<-closed
}
