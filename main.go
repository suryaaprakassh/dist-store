package main

import (
	"dist-store/p2p"
	"log"
)

func main() {
	config := p2p.TcpTransportConfig{
		ListenAddr: ":8080",	
		ShakeHands: p2p.DefaultHandShake,
		Decoder: p2p.DefaultDecoder{},
	}

	tr := p2p.NewTcpTransport(config)

	if err := tr.ListenAndAccept() ; err != nil {
		log.Fatalf("Error starting transport: %v",err.Error())
	}

	select {} 
}
