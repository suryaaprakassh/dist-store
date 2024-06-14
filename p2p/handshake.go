package p2p


// handles the handshake between nodes in a network
type HandShakeFunc func(peer Peer) error

func DefaultHandShake(peer Peer) error {
	return nil
}
