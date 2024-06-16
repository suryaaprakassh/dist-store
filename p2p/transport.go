package p2p

import "net"

// Represents the remote node
type Peer interface {
	RemoteAddr() net.Addr
	Close() error
}

// Handles communication between peers
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
