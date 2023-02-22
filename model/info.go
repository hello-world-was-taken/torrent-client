package model

import (
	// "log"

	// "torrent-dsp/utils"
)

// Info represents the info section of a torrent file
type Info struct {
	PieceLength int64 `bencode:"piece length,omitempty"`
	Pieces      string `bencode:"pieces,omitempty"`
	Private     int64 `bencode:"private,omitempty"`
	Name        string `bencode:"name,omitempty"`
	Length      int64 `bencode:"length,omitempty"`
	Files       []File `bencode:"files,omitempty"`
}


func (i *Info) PiecesToByteArray() [][20]byte {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		// err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes
}