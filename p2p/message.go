package p2p

import "net"

type RPC struct {
	From net.Addr
	Payload []byte
	Peer Peer
}

type Message struct {
	Payload []byte
}