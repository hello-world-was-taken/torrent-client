package utils

import (
	"crypto/sha1"
	"fmt"
	"log"

	"torrent-dsp/model"

	bencode "github.com/zeebo/bencode"
)

func GenerateInfoHash(info model.Info) [20]byte {

	// encode the info back to bencode format
	encodedInfo, err := bencode.EncodeBytes(info)
	if err != nil {
		log.Fatal(err)
	}
	
	// generate sha1 hash from the encoded info
	hash := sha1.Sum(encodedInfo)
	return hash
}