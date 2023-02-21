package client

import (
	// "fmt"
	"log"
	// "net"
	// "time"

	// "torrent-dsp/constant"
	"torrent-dsp/model"
	// "torrent-dsp/utils"
)

func ConnectToTracker() {
	// open torrent file from the current directory and parse it
	// TODO: remove this hardcoded torrent file names
	// ubuntu-22.10-desktop-amd64.iso
	// vlc-media-player
	// 20A4F6FB1C21B5F5D76BAFDA3D64492125F7FAE2
	torrent, err := model.ParseTorrentFile("./torrent-files/debian-11.6.0-amd64-netinst.iso.torrent")
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
	// fmt.Println("Connecting to peers... Length: ", len(peers))
	// ConnectToPeers(peers, torrent)
	// create a client for each peer

	// go SeederMain()  // seed the file

	// peer := model.Peer{}
	// peer.IP = net.IP([]byte{127, 0, 0, 1})
	// peer.Port = 6881
	// peers = []model.Peer{peer}
	// time.Sleep(5 * time.Second)
	StartDownload(torrent, peers)  // leech the file

	// for _, peer := range peers {

	// 	go StartDownload(peer, torrent)
	// 	if err != nil {
	// 		log.Fatal("Error creating client for peer: ", peer.String())
	// 	}

	// 	// fmt.Println("------------  Bit Field  -------------> ", client.BitField)
	// 	// fmt.Println("------------  Choked State  -------------> ", client.ChokedState)
	// }

}
