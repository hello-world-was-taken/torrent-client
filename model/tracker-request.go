package model

import (
	"net/url"
	"strconv"
)

type TrackerRequest struct {
	Info_hash  string `bencode:"info_hash"`
	Peer_id    string `bencode:"peer_id"`
	Port       int    `bencode:"port"`
	Uploaded   int    `bencode:"uploaded"`
	Downloaded int64    `bencode:"downloaded"`
	Left       int64    `bencode:"left"`
	Event      string `bencode:"event"`
}

func (t *TrackerRequest) Encode() string {
	// encode the tracker request
	encodedRequest := url.Values{}
	encodedRequest.Add("info_hash", t.Info_hash)
	encodedRequest.Add("peer_id", t.Peer_id)
	encodedRequest.Add("port", string(strconv.Itoa(t.Port)))
	encodedRequest.Add("uploaded", string(strconv.Itoa(t.Uploaded)))
	// TODO: int64 to string conversion
	encodedRequest.Add("downloaded", string(strconv.Itoa(int(t.Downloaded))))
	encodedRequest.Add("left", string(strconv.Itoa(int(t.Left))))
	encodedRequest.Add("event", t.Event)

	return encodedRequest.Encode()
}