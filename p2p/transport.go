package p2p

//Peer represents a node in the network
type Peer interface {
	Send([]byte) error
}

//Transport is anything that handles the communication between nodes in the network
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}