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