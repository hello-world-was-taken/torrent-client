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

func ClientFactory(peer model.Peer, torrent model.Torrent) (*model.Client, error) {
	client, err := createClient(peer, torrent)
	if err != nil {
		fmt.Println("Error creating client for peer: ", peer.String())
		// log.Fatal(err)
		return nil, nil
	}

	return client, nil
}

// to create a client we need to:
// 1. connect to the peer
// 2. shake hands with the peer
// 3. receive bit field message from the peer
func createClient(peer model.Peer, torrent model.Torrent) (*model.Client, error) {
	conn, err := connectToPeer(peer, torrent)
	if err != nil {
		return nil, err
	}

	// shake hands with the peer
	err = ShakeHandWithPeer(torrent, peer, constant.CLIENT_ID, conn)
	if err != nil {
		return nil, err
	}

	// receive bit field message from the peer
	fmt.Println("Getting Bit Field...")
	bitFieldMessage, err := ReceiveBitFieldMessage(conn)
	fmt.Println("Received Bit field.")
	if err != nil {
		log.Fatal(err)
	}

	// create a new client
	client := &model.Client{
		Peer:        peer,
		BitField:    bitFieldMessage.Payload,
		Conn:        conn,
		ChokedState: 0,
	}

	return client, nil
}

// connect to peer. If not possible then return error
func connectToPeer(peer model.Peer, torrent model.Torrent) (net.Conn, error) {

	conn, err := net.DialTimeout("tcp", peer.String(), constant.CONNECTION_TIMEOUT)
	// TODO: close connection incase of error
	if err != nil {
		fmt.Println("Error connecting to peer: ", peer.String())
		// log.Fatal(err)
		return nil, err
	}

	return conn, nil
}

func ShakeHandWithPeer(torrent model.Torrent, peer model.Peer, clientID string, conn net.Conn) error {
	// create a new connection to the peer
	// conn, err := net.Dial("tcp", peer.String())
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
		fmt.Println("Error sending handshake request to peer: ", peer.String())
		return err
	}

	// check that the infohash in the response matches the infohash of the torrent
	if !bytes.Equal(handshakeResponse.InfoHash[:], torrent.InfoHash[:]) {
		fmt.Println("handshake response infohash does not match torrent infohash")
		return err
	}

	// check that the peer id in the response is different from ours
	if bytes.Equal(handshakeResponse.PeerID[:], utils.ConvertStringToByteArray(constant.CLIENT_ID)[:]) {
		fmt.Println("handshake response peer id matches our peer id")
		return err
	}

	fmt.Println("Handshake successful")
	return nil
}

func ReceiveBitFieldMessage(conn net.Conn) (*model.Message, error) {
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	// send the bitField message
	bitFieldMessageRequest := model.Message{
		MessageID: constant.BIT_FIELD,
		Length:    0,
		Payload:   []byte{0, 0, 0, 0, 0, 0, 0, 0},
	}
	conn.Write(bitFieldMessageRequest.Serialize())

	// receive the bitField message
	bitFieldMessageResponse := model.DeserializeMessage(conn)

	// check that the message is a bit field message
	if bitFieldMessageResponse.MessageID != constant.BIT_FIELD {
		fmt.Println("Expected bit field message")
		// log.Fatal("expected bit field message")
	}

	return bitFieldMessageResponse, nil
}