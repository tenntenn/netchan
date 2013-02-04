package netchan

import (
	"code.google.com/p/go.net/websocket"
)

// Netchan connection to clients or server.
type Conn interface {
	// Connect ch to server or clients.
	Connect(ch interface{}, name []byte) (errCh <-chan error)
	// Get codec
	Codec() websocket.Codec
	// Set codec
	SetCodec(codec websocket.Codec)
}
