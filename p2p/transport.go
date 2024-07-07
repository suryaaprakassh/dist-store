package p2p

import "net"

const (
	StreamingMode=0x1
	NonStreamingMode=0x2
) 

// Represents the remote node
type Peer interface {
	Send([]byte) error
	IsOutbound() bool

	//Sets the peer to file stream mode 
	//Stops the read loop until Stop Stream is called 
	StartStream() 
	
	//Sets the peer to normal reading mode
	//calls the wg.Done() internally
	StopStream() 

	net.Conn
}

// Handles communication between peers
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
