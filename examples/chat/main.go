package main

import (
	"flag"
	"fmt"
	"github.com/tenntenn/netchan"
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

func recieve(closed chan bool) {

	ch := make(chan string)
	fmt.Printf("recieve: %sto%s\n", *remoteHost, *ownHost)
	errCh := netchan.Dial(&ch, nil, *ownHost, *remoteHost, *remoteHost +"to"+*ownHost)

	for {
		select {
		case err := <-errCh:
			if err != nil {
				closed <- true
				return
			}
		case msg := <-ch:
			fmt.Println("Remote say: ", msg)
		case <-closed:
			return
		}
	}
}

func send(closed chan bool) {

	ch := make(chan string)
	done := make(chan bool)
	fmt.Printf("send: %sto%s\n", *ownHost, *remoteHost)
	errCh := netchan.Dial(&ch, nil, *ownHost, *remoteHost, *ownHost + "to" + *remoteHost)
	r := bufio.NewReader(os.Stdin)
	for {
		select {
		case err := <-errCh:
			if err != nil {
				closed <- true
				return
			}
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
			fmt.Println("hoge")
			select {
			case <-done:
				fmt.Println("You said: ", msg)
			case <-time.After(time.Duration(5) * time.Second):
				fmt.Println("time out")
				closed <- true
				return
			}
		}
	}
}

func run() {
	closed := make(chan bool)
	go recieve(closed)
	send(closed)
}

func main() {
	run()
}
