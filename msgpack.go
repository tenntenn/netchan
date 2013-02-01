package netchan

import (
	"code.google.com/p/go.net/websocket"
	msgpack "github.com/ugorji/go-msgpack"
)

// marshal into messagepack.
// This function is used for websocket.Codec.
func msgpackMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = msgpack.Marshal(v, nil)
	return msg, websocket.BinaryFrame, err
}

// unmarshal from messagepack.
// This function is used for websocket.Codec.
func msgpackUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return msgpack.Unmarshal(msg, v, nil)
}
