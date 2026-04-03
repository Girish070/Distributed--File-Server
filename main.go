package main

import (
	"fmt"
	"log"

	"github.com/girish/storage/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAdder: ":3000",
		Handshakefunc: p2p.NOPHandshakeFunc,
		Decode: p2p.LengthPrefixDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listning on port 3000...")

	select{
		
	}
}
