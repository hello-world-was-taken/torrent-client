package client

import (
	"fmt"
	"log"
	"net"

	"torrent-dsp/constant"
	"torrent-dsp/model"
	"torrent-dsp/utils"
)

func ConnectToTracker() {
	// open torrent file from the current directory and parse it
	// TODO: remove this hardcoded torrent file names
	// ubuntu-22.10-desktop-amd64.iso
	// vlc-media-player
	torrent, err := utils.ParseTorrentFile("./torrent-files/ubuntu-22.10-desktop-amd64.iso.torrent")
	if err != nil {
		log.Fatal(err)
	}

	// get a list of peers from the tracker
	peers, err := GetPeersFromTrackers(&torrent)
	if err != nil {
		log.Fatal(err)
	}

	// log the peers
	// utils.LogPeers(peers)

	// connect to the peers
	fmt.Println("Connecting to peers...")
	ConnectToPeers(peers, torrent)

}

// connect to the peers
func ConnectToPeers(peers []model.Peer, torrent model.Torrent) {
	for _, peer := range peers {
		// create a new connection to the peer
		conn, err := net.DialTimeout("tcp", peer.String(), constant.CONNECTION_TIMEOUT)
		// TODO: close connection incase of error
		if err != nil {
			fmt.Println("Error connecting to peer: ", peer.String())
			log.Fatal(err)
		}

		// create a new client with the peer
		client := CreateClient(torrent, peer, constant.CLIENT_ID, conn)
		fmt.Println("------------  Bit Field  -------------> ", client.BitField)

	}
}

func CreateClient(torrent model.Torrent, peer model.Peer, clientID string, conn net.Conn) *model.Client {
	// TODO: use goroutines to connect to multiple peers at the same time
	// shake hands with the peer
	ShakeHandWithPeer(torrent, peer, constant.CLIENT_ID, conn)

	// receive bitfield message from the peer
	fmt.Println("Getting Bit Field...")
	bitFieldMessage, err := ReceiveBitFieldMessage(conn)
	fmt.Println("Received Bit field.")
	if err != nil {
		log.Fatal(err)
	}

	// create a new client
	client := &model.Client{
		Peer:        peer,
		BitField:    bitFieldMessage.Payload,
		Conn:        conn,
		ChokedState: 0,
	}

	return client
}
