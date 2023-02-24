package model

import (
	"encoding/binary"
	"errors"

	// "fmt"
	"log"
	"net"
	"syscall"
	"torrent-dsp/constant"
)

type Client struct {
	Conn        net.Conn
	Peer        Peer
	BitField    []byte
	ChokedState uint8
}

func (client *Client) Interested() {
	msg := Message{MessageID: constant.INTERESTED, Payload: []byte{}}
	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending interested message to peer: %s", err)
	}
}

func (client *Client) Choke() {
	msg := Message{MessageID: constant.CHOKE, Payload: []byte{}}
	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending choke message to peer: %s", err)
	}
}

func (client *Client) UnChoke() {
	msg := Message{MessageID: constant.UN_CHOKE, Payload: []byte{}}
	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending unchoke message to peer: %s", err)
	}

}

// length: 4 bytes, id: 1 byte, index: 4 bytes, begin: 4 bytes, length: 4 bytes -> total 17 bytes
func (client *Client) Request(index uint32, begin uint32, length uint32) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], index)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	binary.BigEndian.PutUint32(payload[8:12], length)
	msg := Message{MessageID: constant.REQUEST, Payload: payload}

	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			// fmt.Println("Peer disconnected")
		}
		// fmt.Printf("Error sending request message to peer: %s", err)
		return err
	}

	return nil
}

func (client *Client) Have(index uint32) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload[0:4], index)
	msg := Message{MessageID: constant.HAVE, Payload: payload}

	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending have message to peer: %s", err)
	}
}
