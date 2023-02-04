// dial to localhost:6889 and send a message
package client

import (
	"log"
	
	"torrent-dsp/utils"
)


func ConnectToTracker() {
	// open torrent file from the current directory and parse it
	torrent, err := utils.ParseTorrentFile("./torrent-files/ubuntu-22.10-desktop-amd64.iso.torrent")
	if err != nil {
		log.Fatal(err)
	}

	// get a list of peers from the tracker
	_, err = GetPeersFromTrackers(&torrent)

}