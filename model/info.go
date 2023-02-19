package model

import (
	"log"

	"torrent-dsp/utils"
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
	// convert the pieces string to a byte array of hashes (20 bytes each)
	pieces := []byte(i.Pieces)
	piecesInByteArray := [][20]byte{}
	
	if len(pieces) % 20 != 0 {
		log.Fatal("Error: pieces length is not a multiple of 20")
		return nil
	}

	for i := 0; i < len(pieces); i += 20 {
		curSlice := [20]byte{}
		end := utils.CalcMin(i+20, len(pieces))
		copy(curSlice[:], pieces[i:end])

		piecesInByteArray = append(piecesInByteArray, curSlice)
	}

	// numHashes := len(pieces) / 20
	// hashes := make([][20]byte, numHashes)

	// for i := 0; i < numHashes; i++ {
	// 	copy(hashes[i][:], pieces[i*20:(i+1)*20])
	// }
	// return hashes

	return piecesInByteArray
}