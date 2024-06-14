package main

import (
	"dist-store/p2p"
	"fmt"
	"log"
)

func main() {
	config := p2p.TcpTransportConfig{
		ListenAddr: ":8080",	
		ShakeHands: p2p.DefaultHandShake,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: func(p p2p.Peer) error {
			return fmt.Errorf("Failed on Peer")
		},
	}

	tr := p2p.NewTcpTransport(config)
	
	go func() {
		for {
			msg := <- tr.Consume()
			log.Printf("message: %v",msg)
		}
	}()

	if err := tr.ListenAndAccept() ; err != nil {
		log.Fatalf("Error starting transport: %v",err.Error())
	}

	select {} 
}
