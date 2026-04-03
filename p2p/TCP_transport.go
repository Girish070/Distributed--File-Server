package p2p

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	//Conn is the underlying connection of the peer
	conn net.Conn
	//is we dial and retrive a conn => outbound == true
	//if we accept and retrive a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOps struct {
	ListenAdder   string
	Handshakefunc HandshakeFunc
	Decode        Decoder
}

type TCPTransport struct {
	TCPTransportOps
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
		peers: make(map[net.Addr]Peer),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAdder)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP connection error: %s\n", err)

		}
		fmt.Printf("New incoming connection: %s\n", conn)

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {

	peer := NewTCPPeer(conn, true)

	if err := t.Handshakefunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	t.mu.Lock()
	t.peers[conn.RemoteAddr()] = peer
	t.mu.Unlock()

	msg := &Message{}
	for {
		if err := t.Decode.Decode(conn, msg); err != nil {
			// ⚠️ Check if the error is just a normal disconnect
			if err == io.EOF {
				fmt.Printf("Client %s disconnected cleanly.\n", conn.RemoteAddr())
				break // Exit the for-loop
			}
			
			// If it's a different error, print it
			fmt.Printf("TCP read error: %s\n", err)
			break // Exit the loop on errors so the goroutine doesn't spin forever
		}
		
		fmt.Printf("Received Message: %s\n", string(msg.Payload))
	}
	
	// Clean up when the loop ends
	conn.Close()
	fmt.Println("Connection closed and goroutine finished.")
}

func (p *TCPPeer) Send(b []byte) error {
	//1. Create 4 byte buffer for the length
	lengthBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBuf, uint32(len(b)))

	//2. Write the length followed by the actual payload
	//Note: In production, I want to do this in a single write call
	//Using a buffered writer or append() to sending two separate TCP packets
	msg := append(lengthBuf, b...)
	_, err := p.conn.Write(msg)
	return err
}
