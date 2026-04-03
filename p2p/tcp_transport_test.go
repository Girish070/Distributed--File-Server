package p2p

import (
	"bytes"
	"net"
	"testing"
)

// Test 1: Verify the server can initialize and listen on a port
func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOps{
		ListenAdder:   ":4000", // Use a different port than main.go to avoid conflicts
		Handshakefunc: NOPHandshakeFunc,
		Decode:        LengthPrefixDecoder{},
	}
	tr := NewTCPTransport(opts)

	err := tr.ListenAndAccept()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// Test 2: Verify our length-prefixed sending and decoding works perfectly
func TestTCPPeerSendAndDecode(t *testing.T) {
	// net.Pipe creates a direct, in-memory connection.
	// Whatever is written to clientConn can be read from serverConn!
	clientConn, serverConn := net.Pipe()

	// 1. Setup the Sender (Client)
	peer := NewTCPPeer(clientConn, true)
	expectedPayload := []byte("Testing the length prefixed protocol!")

	// We run the send in a goroutine because writing to a net.Pipe blocks
	// until the other side reads from it.
	go func() {
		err := peer.Send(expectedPayload)
		if err != nil {
			t.Errorf("failed to send message: %v", err)
		}
		clientConn.Close() // Close connection after sending
	}()

	// 2. Setup the Receiver (Server Decoder)
	decoder := LengthPrefixDecoder{}
	msg := &Message{}

	// Read from the server side of the pipe
	err := decoder.Decode(serverConn, msg)
	if err != nil {
		t.Fatalf("failed to decode message: %v", err)
	}

	// 3. Verify the results
	if !bytes.Equal(expectedPayload, msg.Payload) {
		t.Errorf("Expected payload %q, but got %q", expectedPayload, msg.Payload)
	}
}