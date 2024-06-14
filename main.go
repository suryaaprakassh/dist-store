package main

import (
	"dist-store/p2p"
	"log"
)

func main()  {
	tr := p2p.NewTcpTransport(":8080")

	if err := tr.ListenAndAccept() ; err != nil {
		log.Fatalf("Failed to start transport!")
	}
	
	select{}
}
