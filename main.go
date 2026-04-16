package main

import (
	"log"
	"time"

	"github.com/girish/storage/p2p"
	"github.com/girish/storage/store"
)

func makeServer(listenAddr string, storageroot string, bootstrapNodes []string) *FileServer {
	tcpOpts := p2p.TCPTransportOps{
		ListenAdder:   listenAddr,
		Handshakefunc: p2p.NOPHandshakeFunc,
		Decode:        p2p.LengthPrefixDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	storeOpts := store.StoreOpts{
		Root:              storageroot,
		PathTransformFunc: store.CASPathTransfromFunc,
	}
	localStore := store.NewStore(storeOpts)

	serverOpts := FileServerOpts{
		StorageRoot:    storageroot,
		Transport:      tr,
		Store:          localStore,
		BootStrapNoeds: bootstrapNodes,
	}
	return NewFileServer(serverOpts)
}

func main() {
	// 1. Create Node 1 (The "Seed" Node)
	// It listens on :3000 and has no bootstrap nodes because it is the first one.
	node1 := makeServer(":3000", "node1_data", nil)
	
	// Start Node 1 in the background
	go func() {
		log.Fatal(node1.Start())
	}()

	// Give Node 1 a second to fully boot up its TCP listener
	time.Sleep(1 * time.Second)

	// 2. Create Node 2
	// It listens on :4000, and we tell it to dial Node 1 when it starts!
	node2 := makeServer(":4000", "node2_data", []string{":3000"})
	
	// Start Node 2 (this will block and keep the program running)
	log.Fatal(node2.Start())
}