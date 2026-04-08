package main

import (
	"log"

	"github.com/girish/storage/p2p"
	"github.com/girish/storage/store"
)

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAdder:   ":3000",
		Handshakefunc: p2p.NOPHandshakeFunc,
		Decode:        p2p.LengthPrefixDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	storeOpts := store.StoreOpts{
		Root:              "my_network_data",
		PathTransformFunc: store.CASPathTransfromFunc,
	}
	localStore := store.NewStore(storeOpts)

	serverOpts := FileServerOpts{
		StorageRoot: "my_network_data", // This is where files will be saved later
		Transport:   tr,
		Store:       localStore,
	}
	server := NewFileServer(serverOpts)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
