package model

// File represents a file in the info section of a torrent file
type File struct {
	Length int64 `bencode:"length"`
	Path   []string `bencode:"path"`
}