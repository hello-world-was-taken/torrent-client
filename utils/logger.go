package utils

import (
	"fmt"
	"log"

	"torrent-dsp/model"
)

// print torrent file to console
func LogTorrent(torrent model.Torrent) {
	log.Println("Announce: ", torrent.Announce)
	log.Println("Announce List: ", torrent.AnnounceList)
	log.Println("Comment: ", torrent.Comment)
	log.Println("Created By: ", torrent.CreatedBy)
	log.Println("Creation Date: ", torrent.CreationDate)
	log.Println("Piece Length: ", torrent.Info.PieceLength)
	log.Println("Name: ", torrent.Info.Name)
	log.Println("Length: ", torrent.Info.Length)
}


// print tracker list of peers of tracker response to console
func LogPeers(response *model.TrackerResponse) {
	fmt.Println("Peers: ", response)
}