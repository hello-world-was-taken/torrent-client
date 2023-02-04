package model

type TrackerRequest struct {
	Info_hash  string `bencode:"info_hash"`
	Peer_id    string `bencode:"peer_id"`
	Port       int    `bencode:"port"`
	Uploaded   int    `bencode:"uploaded"`
	Downloaded int64    `bencode:"downloaded"`
	Left       int64    `bencode:"left"`
	Event      string `bencode:"event"`
}