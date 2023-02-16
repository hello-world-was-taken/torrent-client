package model

import (
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
	
	for i := 0; i < len(pieces); i += 20 {
		curSlice := [20]byte{}
		end := utils.CalcMin(i+20, len(pieces))
		copy(curSlice[:], pieces[i:end])

		piecesInByteArray = append(piecesInByteArray, curSlice)
	}

	return piecesInByteArray
}