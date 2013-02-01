netchan
==============

An implementation of netchan using websocket and messagepack.

API document is [here](http://godoc.org/github.com/tenntenn/netchan).

Send
--------

```go
ch := make(chan int)
done := make(chan bool)
select {
	case err := <-netchan.Dial(&ch, done, ":9090", ":8080", "Test"):
		fmt.Println(err.Error())
	case ch <- 100:
		<-done
}
```

Recieve
--------

```go
ch := make(chan int)
select {
	case err := <-netchan.Dial(&ch, nil, ":8080", ":9090", "Test"):
		fmt.Println(err.Error())
	case n := <-ch:
		fmt.Println(n)
}
```
