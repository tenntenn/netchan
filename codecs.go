package netchan

import (
	"code.google.com/p/go.net/websocket"
	"encoding/xml"
	msgpack "github.com/ugorji/go-msgpack"
)

var (
	// Codec of messagepack.
	MsgpackCodec = websocket.Codec{msgpackMarshal, msgpackUnmarshal}
	// Codec of xml.
	XmlCodec = websocket.Codec{xmlMarshal, xmlUnmarshal}
)

// marshal into messagepack.
// This function is used for websocket.Codec.
func msgpackMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = msgpack.Marshal(v)
	return msg, websocket.BinaryFrame, err
}

// unmarshal from messagepack.
// This function is used for websocket.Codec.
func msgpackUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return msgpack.Unmarshal(msg, v, nil)
}

// marshal into xml.
// This function is used for websocket.Codec.
func xmlMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = xml.Marshal(v)
	return msg, websocket.TextFrame, err
}

// unmarshal from xml.
// This function is used for websocket.Codec.
func xmlUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return xml.Unmarshal(msg, v)
}
