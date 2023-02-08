package model

import (
	"fmt"
	"net"
)

type HandShake struct {
	Pstr     string   `bencode:"pstr"`
	InfoHash [20]byte `bencode:"info_hash"`
	PeerID   [20]byte `bencode:"peer_id"` 
}

// serialize the handshake request
// <pstrlen 1 byte><pstr 19 bytes i.e Bit Torrent><reserved 8 byte><info_hash><peer_id 20 bytes>
func (handShake *HandShake) Serialize() []byte {
	// create a buffer to store the handshake request. For bittorrent, request length is 49 + len(pstr) ==> 49 + 19 = 68
	buffer := make([]byte, 49+len(handShake.Pstr))
	buffer[0] = byte(len(handShake.Pstr))
	copy(buffer[1:], []byte(handShake.Pstr))
	copy(buffer[1+len(handShake.Pstr):], make([]byte, 8))
	copy(buffer[1+len(handShake.Pstr)+8:], handShake.InfoHash[:])
	copy(buffer[1+len(handShake.Pstr)+8+20:], handShake.PeerID[:])

	return buffer
}

// deserialize the handshake response byte array into a HandShake struct
func DeserializeHandShake(buffer []byte) *HandShake {
	handShake := &HandShake{}
	pstrLength := int(buffer[0])
	handShake.Pstr = string(buffer[1 : pstrLength+1])
	// TODO: remove magic numbers
	copy(handShake.InfoHash[:], buffer[28:48])
	copy(handShake.PeerID[:], buffer[48:68])

	return handShake
}

func (h *HandShake) Send(conn net.Conn) (*HandShake, error) {
	// serialize the handshake request
	buffer := h.Serialize()

	// send the handshake request
	_, err := conn.Write(buffer)
	if err != nil {
		return &HandShake{}, err
	}

	// read the handshake response
	// handshake response is 68 bytes long <length of Protocol> + <Protocol> + <Reserved> + <InfoHash> + <PeerID>
	buffer = make([]byte, 68)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading handshake response")
		return &HandShake{}, err
	}

	// deserialize the handshake response
	handShake := DeserializeHandShake(buffer)

	fmt.Println("Handshake sent successfully")
	return handShake, nil
}
