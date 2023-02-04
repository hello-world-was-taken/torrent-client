package utils

import (
	"fmt"
	"encoding/binary"
	"errors"
	"net"
	"torrent-dsp/model"
)

// parse peers from a byte array of peers
func PeerParser(peerByte []byte) ([]model.Peer, error) {
	// each peer is 6 bytes long
	numOfPeers := len(peerByte) / 6
	peers := make([]model.Peer, numOfPeers)

	if len(peerByte)%6 != 0 {
		fmt.Println("---------------  size error ---------------", len(peerByte))
		return []model.Peer{}, errors.New("invalid peer byte array")
	}

	// peers is a list of bytes. But from that list of bytes, we need to get the ip address and port
	// of each peer. And the way they are organized is that each peer's IP and port is presented as
	// 6 bytes. The first 4 bytes are the IP address and the last 2 bytes are the port.
	for i := 0; i < numOfPeers; i++ {
		// get the ip address
		ip := net.IP(peerByte[i*6: i*6+4])
		// get the port
		port := binary.BigEndian.Uint16([]byte(peerByte[i*6+4: i*6+6]))
		peers[i] = model.Peer{IP: ip, Port: port}
	}

	return peers, nil
}