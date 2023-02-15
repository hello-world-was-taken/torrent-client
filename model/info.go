package model

// Info represents the info section of a torrent file
type Info struct {
	PieceLength int64 `bencode:"piece length,omitempty"`
	Pieces      string `bencode:"pieces,omitempty"`
	Private     int64 `bencode:"private,omitempty"`
	Name        string `bencode:"name,omitempty"`
	Length      int64 `bencode:"length,omitempty"`
	Files       []File `bencode:"files,omitempty"`
}


func (i *Info) PiecesToByteArray() []byte {
	// convert the pieces string to a byte array of hashes (20 bytes each)
	pieces := []byte(i.Pieces)
	piecesInByteArray := make([]byte, len(pieces)/20*20)
	for i := 0; i < len(pieces); i += 20 {
		copy(piecesInByteArray[i/20*20:], pieces[i:i+20])
	}

	return piecesInByteArray
}