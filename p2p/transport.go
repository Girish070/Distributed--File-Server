package p2p

//Peer represents a node in the network
type Peer interface {}

//Transport is anything that handles the communication 
type Transport interface {
	ListenAndAccept() error
}