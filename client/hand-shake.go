package client

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"

	"torrent-dsp/constant"
	"torrent-dsp/model"
	"torrent-dsp/utils"
)

// TODO: rename file to reflect both handshake and bitfield message
func ShakeHandWithPeer(torrent model.Torrent, peer model.Peer, clientID string, conn net.Conn) {
	// create a new connection to the peer
	conn, err := net.Dial("tcp", peer.String())
	if err != nil {
		log.Fatal(err)
	}

	// convert the client id to a byte array
	clientIDByte := [20]byte{}
	copy(clientIDByte[:], []byte(clientID))

	// create a handshake request
	handshakeRequest := model.HandShake{
		Pstr:     "BitTorrent protocol",
		InfoHash: torrent.InfoHash,
		PeerID:   clientIDByte,
	}

	// send the handshake request
	fmt.Println("Sending handshake request to peer: ", peer.String())
	handshakeResponse, err := handshakeRequest.Send(conn)
	if err != nil {
		log.Fatal(err)
	}

	// check that the infohash in the response matches the infohash of the torrent
	if !bytes.Equal(handshakeResponse.InfoHash[:], torrent.InfoHash[:]) {
		log.Fatal("handshake response infohash does not match torrent infohash")
	}

	// check that the peer id in the response is different from ours
	if bytes.Equal(handshakeResponse.PeerID[:], utils.ConvertStringToByteArray(constant.CLIENT_ID)[:]) {
		log.Fatal("handshake response peer id matches our peer id")
	}

	fmt.Println("Handshake successful")
}

func ReceiveBitFieldMessage(conn net.Conn) (*model.Message, error) {
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	// receive the bitField message
	bitFieldMessage := model.DeserializeMessage(conn)

	// check that the message is a bit field message
	if bitFieldMessage.MessageID != constant.BITFIELD {
		log.Fatal("expected bit field message")
	}

	return bitFieldMessage, nil
}
