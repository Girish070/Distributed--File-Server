package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"sync"

	"github.com/girish/storage/p2p"
	"github.com/girish/storage/store"
)

type FileServerOpts struct {
	ListenAddr     string
	StorageRoot    string //StorageRoot is the folder on hard drive where this node will save files
	Transport      p2p.Transport
	Store          *store.Store
	BootStrapNoeds []string
}

// DataMessage is wire protocol payload!
type MessagePayload struct {
	Command string
	Key     string
	Data    []byte
}

type FileServer struct {
	FileServerOpts
	peerLock sync.Mutex          //Thread sefty for concurrent connectons
	peers    map[string]p2p.Peer //Maps addresses to active network connection
	Ring     *p2p.HashRing       //The decentralized routing table
	quitCh   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	fs := &FileServer{
		FileServerOpts: opts,
		peers:          make(map[string]p2p.Peer),
		Ring:           p2p.NewHashring(),
		quitCh:         make(chan struct{}),
	}
	fs.Ring.AddNode(opts.ListenAddr)

	return fs
}

// Start boots up the transport layer and begins processing messages
func (s *FileServer) Start() error {
	fmt.Printf("Starting File Server... (Storage Root: %s)\n", s.StorageRoot)

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	// Bootstrap the network! Connect to all known peers.
	s.bootstrapNetwork()

	s.loop()
	return nil
}

// BootStrapNetwork iterates through all provideed nodes and attempts to connect
func (s *FileServer) bootstrapNetwork() {
	for _, addr := range s.BootStrapNoeds {
		if len(addr) == 0 {
			continue
		}
		fmt.Printf("Attempting to connect with bootstrap node at %s...\n", addr)
		// running this in goroutine so one slow connection dosen't block the others
		go func(peerAddr string) {
			if err := s.Transport.Dial(peerAddr); err != nil {
				fmt.Printf("Failed to dial bootstrap node %s: %s\n", peerAddr, err)
			}
		}(addr)
	}
}

func (s *FileServer) loop() {
	for {
		select {
		case rpc := <-s.Transport.Consume():
			s.handleMessage(rpc)
		case <-s.quitCh:
			fmt.Println("File Server shutting down")
			return
		}
	}
}

func (s *FileServer) handleMessage(rpc p2p.RPC) {
	s.registerPeer(rpc.Peer, rpc.From.String())

	var msg MessagePayload
	if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
		fmt.Printf("Failed to decode network payload: %s\n", err)
		return
	}

	switch msg.Command {
	case "PUT":
		s.handlePutCommand(rpc.From.String(), msg)
	case "GET":
		s.handleGetCommand(rpc)
	case "DELETE":
		s.handleDeleteCommand(rpc.From.String(), msg)
	default:
		fmt.Printf("Unknown command recived: %s\n", msg.Command)
	}
}

func (s *FileServer) handlePutCommand(from string, msg MessagePayload) {
	fmt.Printf("Receiving 'PUT' command for '%s' from %s\n", msg.Key, from)

	err := s.Store.WriteStream(msg.Key, bytes.NewReader(msg.Data))
	if err != nil {
		fmt.Printf("Error storing file: %s\n", err)
		return
	}
	fmt.Println("File successfully saved to disk!")

	err = s.broadcast(msg)
	if err != nil {
		fmt.Printf("Error broadcasting file: %s\n", err)
	}
}

func (s *FileServer) handleGetCommand(rpc p2p.RPC) {
	var msg MessagePayload
	gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg)

	fmt.Printf("Receiving 'GET' command for '%s' from %s\n", msg.Key, rpc.From)

	//1 Open the file stream
	r, err := s.Store.ReadStream(msg.Key)
	if err != nil {
		fmt.Printf("Error reading file from disk: %s\n", err)
		return
	}
	defer r.Close()

	//2 Read file into memory
	fileBytes, err := io.ReadAll(r)
	if err != nil {
		fmt.Printf("Error reading file bytes: %s\n", err)
		return
	}

	//3 Send it back to the client
	if err := rpc.Peer.Send(fileBytes); err != nil {
		fmt.Printf("Error sending file to peer: %s\n", err)
	}
	fmt.Printf("Successfully sent file '%s' back to client!\n", msg.Key)
}

// handleDeleteCommand handles request to remove files from the network
func (s *FileServer) handleDeleteCommand(from string, msg MessagePayload) {
	fmt.Printf("Receiving 'DELETE' command for '%s' from %s\n", msg.Key, from)

	err := s.Store.Delete(msg.Key)
	if err != nil {
		fmt.Printf("Error deleting file from disk: %s\n", err)
		return
	}
	fmt.Printf("File '%s' and its directories successfully deleted!\n", msg.Key)
}

func (s *FileServer) registerPeer(peer p2p.Peer, addr string) {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	if _, exists := s.peers[addr]; !exists {
		s.peers[addr] = peer
		s.Ring.AddNode(addr)
		fmt.Printf("New peer Added %s\n", addr)
	}
}

func (s *FileServer) broadcast(msg MessagePayload) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	payloadBuff := new(bytes.Buffer)
	if err := gob.NewEncoder(payloadBuff).Encode(msg); err != nil {
		return err
	}
	payloadBytes := payloadBuff.Bytes()

	for addr, peer := range s.peers {
		fmt.Printf("📡 Broadcasting file '%s' to peer: %s\n", msg.Key, addr)
		if err := peer.Send(payloadBytes); err != nil {
			fmt.Printf("Failed to broadcast to peer %s: %s\n", addr, err)
		}
	}
	return nil
}
