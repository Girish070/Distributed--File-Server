package main

import (
	"fmt"
	"log"

	"github.com/girish/storage/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listning on port 3000...")

	select{
		
	}
}
