package main

import (
	"log"

	"github.com/girish/storage/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAdder: ":3000",
		Handshakefunc: p2p.NOPHandshakeFunc,
		Decode: p2p.LengthPrefixDecoder{},
	}
	tr:= p2p.NewTCPTransport(tcpOpts)

	serverOpts := FileServerOpts{
		StorageRoot: "my_network_data", // This is where files will be saved later
		Transport: tr,
	}
	server := NewFileServer(serverOpts)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
