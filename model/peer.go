package model

import "net"

type Peer struct {
	IP   net.IP `bencode:"ip"`
	Port uint16   `bencode:"port"`
}