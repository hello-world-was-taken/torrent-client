package model

import (
	"fmt"
	"os"
	bencode "github.com/zeebo/bencode"

)

// parse a torrent file and return a Torrent struct
func ParseTorrentFile(filename string) (Torrent, error) {
	// open the file
	file, err := os.Open(filename)
	if err != nil {
		return Torrent{}, err
	}
	defer file.Close()

	// decode the file
	var torrent = Torrent{}
	err = bencode.NewDecoder(file).Decode(&torrent)
	if err != nil {
		fmt.Println("Encountered error while decoding")
		return Torrent{}, err
	}

	// generate the info hash
	torrent.GenerateInfoHash()
	fmt.Println("Torrent file parsed successfully")
	return torrent, nil
}