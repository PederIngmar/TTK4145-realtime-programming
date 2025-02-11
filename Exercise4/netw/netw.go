package main

import (
	"Ex4/netw/network/bcast"
	"Ex4/netw/network/localip"
	"Ex4/netw/network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)


type HelloMsg struct {
	Message string
	Iter    int
}

type AliveMsg struct {
	//lastAlive 	time.Time
	Msg 		string
	Count int
}

func main() {

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	
	//helloTx := make(chan HelloMsg)
	//helloRx := make(chan HelloMsg)

	aliveTx := make(chan AliveMsg)
	aliveRx := make(chan AliveMsg)

	go bcast.Transmitter(16569, aliveTx)
	go bcast.Receiver(16569, aliveRx)

	// The example message. We just send one of these every second.
	go func() {
		aliveMsg := AliveMsg{"Alive", 0}
		for {
			aliveMsg.Count++
			aliveTx <- aliveMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-aliveRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}
