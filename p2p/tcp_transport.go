package p2p

import (
	"log/slog"
	"net"
	"sync"
)

// Represents a node in tcp network
type TcpPeer struct {
	conn net.Conn
	// If dialed its outBound else its inbound conn
	outBound bool
}

// handles the peer cleanUp logic
func (p *TcpPeer) Close() error {
	return p.conn.Close()
}

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		conn:     conn,
		outBound: outbound,
	}
}

type TcpTransportConfig struct {
	ListenAddr string
	ShakeHands HandShakeFunc
	Decoder    Decoder
}

type TcpTransport struct {
	TcpTransportConfig
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTcpTransport(config TcpTransportConfig) *TcpTransport {
	return &TcpTransport{
		TcpTransportConfig: config,
	}
}

func (t *TcpTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil
}

func (t *TcpTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			slog.Error("Tcp accept error", err)
		}

		go t.handleConn(conn)
	}
}

func (t *TcpTransport) handleConn(conn net.Conn) {
	peer := NewTcpPeer(conn, true)

	if err := t.ShakeHands(peer); err != nil {
		slog.Error("Error in Handshake", err)
		peer.Close()
		return
	}
	slog.Info("New Connection", "Addr", peer.conn.LocalAddr().String())

	var msg Message
	msg.From = conn.RemoteAddr()

	//message read loop
	for {
		if err := t.Decoder.Decode(peer.conn, &msg); err != nil {
			slog.Error("Error in Decoding", err)

			//TODO: fix abrubt connection close
			peer.Close()
			return
		}
		slog.Info("New", "Message", string(msg.Payload))
	}
}
