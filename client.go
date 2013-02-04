package netchan

import (
	"code.google.com/p/go.net/websocket"
	"encoding/base64"
	"fmt"
	"net"
)

// client connection.
type clientConn struct {
	origin  string
	urlBase string
	codec   websocket.Codec
}

func Dial(addr string) (Conn, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	host := "localhost"
	if tcpAddr.IP != nil {
		host = tcpAddr.IP.String()
	}

	port := tcpAddr.Port
	urlBase := fmt.Sprintf("ws://%s:%d", host, port)
	origin := fmt.Sprintf("http://%s:%d", host, port)

	return &clientConn{origin, urlBase, MsgpackCodec}, nil
}

// Importing channel from websocket sever.
func (cnn *clientConn) Connect(ch interface{}, name []byte) (errCh <-chan error) {
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

	encodedName := base64.URLEncoding.EncodeToString(name)
	url := fmt.Sprintf("%s/%s", cnn.urlBase, encodedName)
	ws, err := websocket.Dial(url, "", cnn.origin)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if v, ok := chpv.Recv(); ok {
				cnn.codec.Send(ws, v)
			}
		}
	}()

	return
}

// get codec
func (cnn *clientConn) Codec() websocket.Codec {
	return cnn.codec
}

// set codec
func (cnn *clientConn) SetCodec(codec websocket.Codec) {
	cnn.codec = codec
}
