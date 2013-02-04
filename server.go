package netchan

import (
	"code.google.com/p/go.net/websocket"
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
)

// server connection.
type serverConn struct {
	httpServer *http.Server
	handlers   *http.ServeMux
	channels   map[string]bool
	codec      websocket.Codec
}

// Serve websocket connection.
func Serve(addr string) Conn {

	// create http server for handshaking.
	handlers := http.NewServeMux()
	server := &http.Server{Addr: addr, Handler: handlers}

	// start http server
	go server.ListenAndServe()

	return &serverConn{server, handlers, make(map[string]bool), MsgpackCodec}
}

// Exporting given channel by a websocket handler.
func (cnn *serverConn) Connect(ch interface{}, name []byte) (errCh <-chan error) {
	errSender := make(chan error)
	errCh = errSender // read only
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			default:
				errSender <- fmt.Errorf("Error: %v", r)
			case error:
				err := r.(error)
				errSender <- err
			}
		}
	}()

	chpv := valueChan(ch)

	// regist handler
	cnn.handle(chpv, name)

	return
}

func (cnn *serverConn) handle(chpv reflect.Value, name []byte) {

	handler := func(ws *websocket.Conn) {

		// send
		go func() {
			for {
				if v, ok := chpv.Recv(); ok {
					cnn.codec.Send(ws, v)
				} else {
					panic(Closed)
				}
			}
		}()

		// recieve
		for {
			v := reflect.New(chpv.Type().Elem())
			cnn.codec.Receive(ws, v.Interface())
			chpv.Send(reflect.Indirect(v))
		}
	}

	encodedName := base64.URLEncoding.EncodeToString(name)
	if _, ok := cnn.channels[encodedName]; !ok {
		cnn.channels[encodedName] = true
		pattern := fmt.Sprintf("/%s", encodedName)
		cnn.handlers.Handle(pattern, websocket.Handler(handler))
	}
}

// get codec
func (cnn *serverConn) Codec() websocket.Codec {
	return cnn.codec
}

// set codec
func (cnn *serverConn) SetCodec(codec websocket.Codec) {
	cnn.codec = codec
}
