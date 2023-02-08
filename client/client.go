package client

import (
	"log"
	"net"
	"fmt"

	"torrent-dsp/utils"
	"torrent-dsp/model"
	"torrent-dsp/constant"
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
		if err != nil {
			fmt.Println("Error connecting to peer: ", peer.String())
			// log.Fatal(err)
		}
		
		if err == nil {
			// shake hands with the peer
			ShakeHandWithPeer(torrent, peer, constant.CLIENT_ID, conn)
		}
		
	}
}