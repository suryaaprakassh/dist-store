package main

import (
	"dist-store/p2p"
	"log"
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
	s1 := makeServer(":3000", "TestDir", "")
	s2 := makeServer(":4000", "TestDir", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()

	log.Fatal(s2.Start())
}
