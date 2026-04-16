package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/girish/storage/p2p"
	"github.com/girish/storage/store"
)

type FileServerOpts struct {
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
	quitCh chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		quitCh:         make(chan struct{}),
	}
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
