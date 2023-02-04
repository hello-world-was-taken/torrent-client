package model

import (
	"net/url"
	"strconv"
)

type TrackerRequestParams struct {
	Info_hash  [20]byte `bencode:"info_hash"`
	Peer_id    string `bencode:"peer_id"`
	Port       int    `bencode:"port"`
	Uploaded   int    `bencode:"uploaded"`
	Downloaded int64    `bencode:"downloaded"`
	Left       int64    `bencode:"left"`
	Compact    int    `bencode:"compact"`
	Event      string `bencode:"event"`
}

func (t *TrackerRequestParams) Encode() string {
	// encode the tracker request
	encodedRequest := url.Values{
		"info_hash":  []string{string(t.Info_hash[:])},
		"peer_id":    []string{t.Peer_id},
		"port":       []string{strconv.Itoa(t.Port)},
		"uploaded":   []string{strconv.Itoa(t.Uploaded)},
		"compact":    []string{strconv.Itoa(t.Compact)},
		// TODO: int64 to string conversion
		"downloaded": []string{strconv.Itoa(int(t.Downloaded))},
		"left":       []string{strconv.Itoa(int(t.Left))},
		"event":      []string{t.Event},
	}

	return encodedRequest.Encode()
}