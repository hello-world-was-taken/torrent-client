// Package torrentparser provides a simple interface for parsing torrent files
package utils

import (
	"fmt"
	"os"
	bencode "github.com/zeebo/bencode"
	"torrent-dsp/model"
)

// parse a torrent file and return a Torrent struct
func ParseTorrentFile(filename string) (model.Torrent, error) {
	// open the file
	file, err := os.Open(filename)
	if err != nil {
		return model.Torrent{}, err
	}
	defer file.Close()

	// decode the file
	var torrent = model.Torrent{}
	err = bencode.NewDecoder(file).Decode(&torrent)
	if err != nil {
		fmt.Println("Encountered error while decoding")
		return model.Torrent{}, err
	}

	// generate the info hash
	infoHash := GenerateInfoHash(torrent.Info)
	torrent.InfoHash = infoHash
	fmt.Println("Torrent file parsed successfully")
	return torrent, nil
}