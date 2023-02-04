// dial to localhost:6889 and send a message
package client

import (
	"fmt"
	"log"
	// "net"
	"torrent-dsp/utils"

	// "github.com/anacrolix/torrent"
)


func ConnectToTracker() {
	// open torrent file from the current directory and parse it
	torrent, err := utils.ParseTorrentFile("./torrent-files/ubuntu-14.04-desktop-amd64+mac.iso.torrent")
	if err != nil {
		log.Fatal(err)
	}

	// log print torrent file
	// utils.LogTorrent(torrent)
	fmt.Println(">>", torrent.InfoHash)

	// connect to tracker from torrent file
	// conn, err := net.Dial("tcp", torrent.Announce)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()

	// // get a list of peers from the connected tracker
	// peers, err := GetPeersFromTracker(torrent, conn)

}