package p2p

//Peer represents a node in the network
type Peer interface {
	Send([]byte) error
}

//Transport is anything that handles the communication 
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}