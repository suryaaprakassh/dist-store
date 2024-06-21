package main

import (
	"bytes"
	"dist-store/p2p"
	"encoding/gob"
	"io"
	"log/slog"
	"sync"
)

type ServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type Server struct {
	ServerOpts
	store  *Store
	quitch chan struct{}

	peerLock sync.Mutex
	peer     map[string]p2p.Peer
}

func (s *Server) Stop() {
	close(s.quitch)
}

func (s *Server) loop() {
	//TODO: have a defer func to do clean up
	for {
		select {
		case msg := <-s.Transport.Consume():
			var p Payload
			slog.Debug("transport consume", "msg", msg)
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
				slog.Error("Error decoding", err)
				return
			}
			if err := s.handlePayload(&p); err != nil {
				slog.Error("Error handling payload", err)
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *Server) handlePayload(p *Payload) error {
	switch p.Action {
	case Save:
		s.StoreData(p.Key, bytes.NewReader(p.Data))
	case Delete:
	//TODO: delete a file
	default:
		slog.Warn("Unknown Request with ", "action", p.Action)
	}

	return nil
}

func (s *Server) bootStrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		//return if no bootstrap nodes configured
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				slog.Error("Error dialing addr", err)
			}
		}(addr)
	}
	return nil
}

func (s *Server) Broadcast(p *Payload) error {
	peers := []io.Writer{}

	buf := new(bytes.Buffer)

	for _, peer := range s.peer {
		if peer.IsOutbound() {
			peers = append(peers, peer)
		}
	}

	//using the multi writer
	mw := io.MultiWriter(peers...)

	err := gob.NewEncoder(buf).Encode(p)

	if err != nil {
		slog.Error("Error in broadcast", err)
	}

	gob.NewEncoder(mw).Encode(buf.Bytes())

	return nil
}

func (s *Server) StoreData(key string, r io.Reader) error {

	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)
	//store the key to the current node
	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	p := Payload{
		Action: Save,
		Key:    key,
		Data:   buf.Bytes(),
	}

	slog.Debug("created payload", "payload", p)

	//share stuff with other nodes
	return s.Broadcast(&p)
}

func (s *Server) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peer[p.RemoteAddr().String()] = p

	slog.Info("New Peer Added", "Addr", p.RemoteAddr().String())
	return nil
}

// starts the server with the transport specifies in server opts
func (s *Server) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootStrapNetwork()
	//this blocks when the Start is called
	s.loop()
	return nil
}

func NewServer(opts ServerOpts) *Server {
	storeOpts := StorageOpts{
		PathTransformFunc: opts.PathTransformFunc,
		RootDir:           opts.StorageRoot,
	}
	return &Server{
		ServerOpts: opts,
		store:      NewStore(storeOpts),
		quitch:     make(chan struct{}),
		peer:       make(map[string]p2p.Peer),
	}
}

// // for gob initialization
// func init() {
// 	gob.Register(StoreDataPayload{})
// }
