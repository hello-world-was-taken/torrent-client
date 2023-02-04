package model

type Peer struct {
	IP   string `bencode:"ip"`
	Port int   `bencode:"port"`
	Peer_id string `bencode:"peer_id"`
}