package main

import (
	"bytes"
	"dist-store/p2p"
	"log"
	"log/slog"
	"time"
)

func makeServer(listenAddr, root string, nodes ...string) *Server {
	config := p2p.TcpTransportConfig{
		ListenAddr: listenAddr,
		ShakeHands: p2p.DefaultHandShake,
		Decoder:    p2p.DefaultDecoder{},
		//TODO: have a onpeer func
	}

	transport := p2p.NewTcpTransport(config)

	serverOpts := ServerOpts{
		ListenAddr:        listenAddr,
		StorageRoot:       root,
		PathTransformFunc: CASPathTransformFunc,
		Transport:         transport,
		BootstrapNodes:    nodes,
	}

	s := NewServer(serverOpts)

	//TODO: fix poor design implementation
	transport.OnPeer = s.OnPeer

	return s
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	s1 := makeServer(":3000", ":3000_dir", "")
	s2 := makeServer(":4000", ":4000_dir", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(time.Second * 2)
	go s2.Start()
	time.Sleep(time.Second * 2)

	data := bytes.NewReader([]byte("test data"))

	err := s2.StoreData("test", data)
	if err != nil {
		log.Fatalf(err.Error())
	}

	select {}
}
