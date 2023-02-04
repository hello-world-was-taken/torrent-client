package model

// Torrent represents a torrent file
type Torrent struct {
	Announce     string `bencode:"announce,omitempty"`
	AnnounceList [][]string `bencode:"announce-list,omitempty"`
	Comment      string `bencode:"comment,omitempty"`
	CreatedBy    string `bencode:"created by,omitempty"`
	CreationDate int64 `bencode:"creation date,omitempty"`
	// ignore attribute if it is not present
	Encoding string `bencode:"encoding,omitempty"`
	Info         Info `bencode:"info,omitempty"`
	InfoHash    [20]byte
}