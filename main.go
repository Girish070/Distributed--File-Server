package main

import (
	"log"
	"os"
	"strings"

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
		ListenAddr:     listenAddr,
		StorageRoot:    storageroot,
		Transport:      tr,
		Store:          localStore,
		BootStrapNoeds: bootstrapNodes,
	}
	return NewFileServer(serverOpts)
}

func main() {
	// 1. Read the network port from Docker (default to :3000)
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":3000"
	}

	// 2. Read the storage folder from Docker
	storageRoot := os.Getenv("STORAGE_ROOT")
	if storageRoot == "" {
		storageRoot = "network_data"
	}

	// 3. Read the comma-separated list of peers to dial on startup
	bootstrapNodesStr := os.Getenv("BOOTSTRAP_NODES")
	var bootstrapNodes []string
	if bootstrapNodesStr != "" {
		bootstrapNodes = strings.Split(bootstrapNodesStr, ",")
	}

	// 4. Spin up the single node
	server := makeServer(listenAddr, storageRoot, bootstrapNodes)
	log.Fatal(server.Start())
}
