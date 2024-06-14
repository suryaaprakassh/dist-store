package p2p

//Represents the remote node
type Peer interface {
	Close() error
}

//Handles communication between peers 
type Transport interface {
	ListenAndAccept() error
}


