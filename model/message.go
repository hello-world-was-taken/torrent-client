package model

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type Message struct {
	Length uint32    `bencode:"length"`
	MessageID uint8    `bencode:"message_id"`
	Payload []byte `bencode:"payload"`
}

// serialize the message into a byte array
// <msgLength 4 byte> <message id 1 byte> <payload>
func (message *Message) Serialize() []byte {
	msgLength := uint32(len(message.Payload) + 1) // +1 for the message id
	buffer := make([]byte, msgLength + 4) // +4 for the msgLength prefix
	binary.BigEndian.PutUint32(buffer[0:4], msgLength)
	buffer[4] = byte(message.MessageID)
	copy(buffer[5:], message.Payload)

	return buffer
}


// deserialize the message byte array into a Message struct
func DeserializeMessage(conn net.Conn) (*Message, error) {
	length := make([]byte, 4)
	_, err := io.ReadFull(conn, length)
	// _, err := conn.Read(length)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection closed")
			return nil, err
		}
		fmt.Println("Error reading message length: %s", err)
		return nil, err
	}

	msgLength := binary.BigEndian.Uint32(length)

	if msgLength == 0 {
		return nil, nil
	}

	buffer := make([]byte, msgLength)
	_, err = io.ReadFull(conn, buffer)
	// _, err = conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading message: %s", err)
		return nil, err
	}

	// fmt.Println("Message length: --------> ", msgLength)
	message := &Message{}

	if msgLength == 0 {
		// keep alive message
		return message, nil
	}

	if msgLength == 1 {
		// either choke, unchoke, interested, not interested
		message.MessageID = buffer[0]
		return message, nil
	}
	message.Length = msgLength
	message.MessageID = buffer[0]
	message.Payload = buffer[1:]
	// fmt.Println("Successfully parsed the message")
	return message, nil
}