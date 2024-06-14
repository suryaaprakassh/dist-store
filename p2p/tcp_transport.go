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

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		conn:     conn,
		outBound: outbound,
	}
}

type TcpTransport struct {
	listenAddr string
	listener   net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTcpTransport(listenAddr string) *TcpTransport {
	return &TcpTransport{
		listenAddr: listenAddr,
	}
}

func (t *TcpTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop();
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
	slog.Info("New Connection", "Addr", peer.conn.LocalAddr().String())
}
