package utils

import (
	"crypto/sha1"
	"fmt"
	"log"

	"torrent-dsp/model"
	bencode "github.com/zeebo/bencode"
)

func GenerateInfoHash(info model.Info) string {

	// encode the info back to bencode format
	encodedInfo, err := bencode.EncodeBytes(info)
	if err != nil {
		fmt.Print("Encountered error while encoding info")
		log.Fatal(err)
	}
	
	// generate sha1 hash from the encoded info
	hash := sha1.Sum(encodedInfo)
	r := fmt.Sprintf("%x", hash)
	fmt.Println("Hash: ", r)
	return r
}