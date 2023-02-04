// dial to localhost:6889 and send a message
package client

import (
	"log"
	
	"torrent-dsp/utils"
)


func ConnectToTracker() {
	// open torrent file from the current directory and parse it
	// TODO: remove this hardcoded torrent file names
	// ubuntu-22.10-desktop-amd64.iso
	// vlc-media-player
	torrent, err := utils.ParseTorrentFile("./torrent-files/vlc-media-player.torrent")
	if err != nil {
		log.Fatal(err)
	}

	// get a list of peers from the tracker
	_, err = GetPeersFromTrackers(&torrent)

}