package model

import (
	"crypto/sha1"
	"log"

	bencode "github.com/zeebo/bencode"
)

// Torrent represents a torrent file
type Torrent struct {
	Announce     string     `bencode:"announce,omitempty"`
	AnnounceList [][]string `bencode:"announce-list,omitempty"`
	Comment      string     `bencode:"comment,omitempty"`
	CreatedBy    string     `bencode:"created by,omitempty"`
	CreationDate int64      `bencode:"creation date,omitempty"`
	// ignore attribute if it is not present
	Encoding string `bencode:"encoding,omitempty"`
	Info     Info   `bencode:"info,omitempty"`
	InfoHash [20]byte
}

func (t *Torrent) GenerateInfoHash() {
	// encode the info back to bencode format
	encodedInfo, err := bencode.EncodeBytes(t.Info)
	if err != nil {
		log.Fatal(err)
	}

	// generate sha1 hash from the encoded info
	hash := sha1.Sum(encodedInfo)
	t.InfoHash = hash
}


func (t *Torrent) CalculateRange(index int) (begin int, end int) {
	begin = index * int(t.Info.PieceLength)
	end = begin + int(t.Info.PieceLength)
	if end > int(t.Info.Length) {
		end = int(t.Info.Length)
	}
	return begin, end
}


func (t *Torrent) CalcRequestSize(index int) int {
	begin, end := t.CalculateRange(index)
	return end - begin
}
