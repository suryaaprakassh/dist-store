package p2p

import "net"

// Represents the remote node
type Peer interface {
	Send([]byte) error
	net.Conn
}

// Handles communication between peers
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
