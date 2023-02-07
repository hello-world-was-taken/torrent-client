package model

import (
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP `bencode:"ip"`
	Port uint16   `bencode:"port"`
}

func (peer *Peer) String() string {
	return peer.IP.String() + ":" + strconv.Itoa(int(peer.Port))
}