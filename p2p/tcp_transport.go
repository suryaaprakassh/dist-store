package p2p

import (
	"errors"
	"log/slog"
	"net"
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

// handles the peer cleanUp logic
func (p *TcpPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *TcpPeer) Send(buf []byte) error {
	_,err := p.conn.Write(buf)
	return err  
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
	OnPeer     func(Peer) error
}

type TcpTransport struct {
	TcpTransportConfig
	listener net.Listener
	rpcch    chan RPC
}

func NewTcpTransport(config TcpTransportConfig) *TcpTransport {
	return &TcpTransport{
		TcpTransportConfig: config,
		rpcch:              make(chan RPC),
	}
}

func (t *TcpTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	slog.Info("Tcp Transport Listening", "On", t.ListenAddr)
	return nil
}

// implements transport interface
// returns readonly channel of type RPC
func (t *TcpTransport) Consume() <-chan RPC {
	return t.rpcch
}

// implements the transport interface
func (t *TcpTransport) Close() error {
	close(t.rpcch)
	return t.listener.Close()
}

// implements the transport interface
func (t *TcpTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TcpTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			slog.Info("Tcp Transport Closed")
			return
		}
		if err != nil {
			slog.Error("Tcp accept error", err)
		}

		go t.handleConn(conn, false)
	}
}

func (t *TcpTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	peer := NewTcpPeer(conn, outbound)

	//clean up logic
	defer func() {
		slog.Error("Droping peer:", err)
		peer.Close()
	}()

	if err = t.ShakeHands(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// slog.Info("New Connection", "Addr", peer.conn.LocalAddr().String())

	var msg RPC
	msg.From = conn.RemoteAddr()
	//message read loop
	for {
		err = t.Decoder.Decode(peer.conn, &msg)

		//TODO: fix net.OpError
		if err != nil {
			return
		}
		t.rpcch <- msg
	}
}
