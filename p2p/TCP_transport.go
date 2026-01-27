package p2p

import (
	"fmt"
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
		fmt.Printf("TCP handshake error: %s\n, err")
		return
	}

	msg := &Message{}
	// buf := make([]byte, 2000)
	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Printf("TCP Error: %s\n,", err)
		// }
		if err := t.Decode.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
		fmt.Printf("Message: %+v\n", msg)
	}
}
