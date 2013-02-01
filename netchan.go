package netchan

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	msgpack "github.com/ugorji/go-msgpack"
	"net"
	"net/http"
	"net/url"
	"reflect"
)

func msgpackMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = msgpack.Marshal(v, nil)
	return msg, websocket.BinaryFrame, err
}

func msgpackUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return msgpack.Unmarshal(msg, v, nil)
}

var (
	Codec    = websocket.Codec{msgpackMarshal, msgpackUnmarshal}
	Protocol = "http"
)

var (
	Done = fmt.Errorf("Done chan")
)

type connMap struct {
	handlers map[string]string
	serveMux *http.ServeMux
}

func newConnMap() *connMap {
	handlers := make(map[string]string)
	serveMux := http.NewServeMux()
	return &connMap{handlers, serveMux}
}

type NetChan struct {
	connMap_ map[string]*connMap
}

func New() *NetChan {
	connMap_ := make(map[string]*connMap)
	return &NetChan{connMap_}
}

func (self *NetChan) Dial(ch interface{}, done chan<- bool, dst, src, name string) (errCh <-chan error) {
	errChBi := make(chan error)
	errCh = errChBi

	var err error
	defer func() {
		if err != nil {
			go func() {
				errChBi <- err
			}()
		}
	}()

	// assert pointer
	chv := reflect.ValueOf(ch)
	if chv.Kind() != reflect.Ptr {
		err = fmt.Errorf("First argument must be pointer of channel.")
		return
	}

	// assert chan
	chpv := reflect.Indirect(chv)
	if chpv.Kind() != reflect.Chan {
		err = fmt.Errorf("First argument must be pointer of channel.")
		return
	}

	// Start http server
	key := fmt.Sprintf("%s->%s", src, dst)
	if connMap_, ok := self.connMap_[key]; !ok {
		connMap_ = newConnMap()
		self.connMap_[key] = connMap_
		server := &http.Server{Addr: src, Handler: connMap_.serveMux}
		go server.ListenAndServe()
	}

	// regist handler
	connMap_ := self.connMap_[key]
	if _, ok := connMap_.handlers[name]; !ok {

		// recieve
		handler := func(ws *websocket.Conn) {
			for {
				v := reflect.New(chpv.Type().Elem())
				Codec.Receive(ws, v.Interface())
				chpv.Send(reflect.Indirect(v))
				Codec.Send(ws, true)
			}
		}
		path := fmt.Sprintf("/%s", url.QueryEscape(name))
		connMap_.handlers[name] = path
		connMap_.serveMux.Handle(path, websocket.Handler(handler))

		// send
		go func() {
			var err error
			defer func() {
				if err != nil {
					go func() {
						errChBi <- err
					}()
				}
			}()

			tcpAddr, _ := net.ResolveTCPAddr("tcp", dst)
			host := "localhost"
			if tcpAddr.IP != nil {
				host = tcpAddr.IP.String()
			}
			port := tcpAddr.Port
			url := fmt.Sprintf("ws://%s:%d%s", host, port, path)
			origin := fmt.Sprintf("%s://%s:%d", Protocol, host, port)

			var ws *websocket.Conn
			ws, err = websocket.Dial(url, "", origin)
			if err != nil {
				return
			}

			for {
				if v, ok := chpv.Recv(); ok {
					Codec.Send(ws, v)
					var res interface{}
					Codec.Receive(ws, &res)
					if done != nil {
						done <- true
					}
				} else {
					break
				}
			}
		}()
	}

	return
}

// Default netchan mapper.
var defaultNetChan = New()

func Dial(ch interface{}, done chan<- bool, dst, src, name string) (errCh <-chan error) {
	return defaultNetChan.Dial(ch, done, dst, src, name)
}
