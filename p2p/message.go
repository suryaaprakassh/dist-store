package p2p

import "net"

// Represents anything that is transferred
// between nodes in a network
type RPC struct {
	From    net.Addr
	Payload []byte
}
