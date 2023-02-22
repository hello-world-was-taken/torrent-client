package client

import (
    "encoding/binary"
	"encoding/hex"
    "math/rand"
	"fmt"
	"net"
)

func UdpCreateConnectionRequest() ([]byte, int) {
    connection_id := int64(0x41727101980) // default connection id
    action := int32(0) // action (0 = give me a new connection id)
    transaction_id := int(rand.Int31n(255))

    buffer := make([]byte, 0, 16)
    binary.BigEndian.PutUint64(buffer, uint64(connection_id)) // first 8 bytes is connection id
    binary.BigEndian.PutUint32(buffer[8:], uint32(action)) // next 4 bytes is action
    binary.BigEndian.PutUint32(buffer[12:], uint32(transaction_id)) // next 4 bytes is transaction id

    return buffer, int(transaction_id)
}


func UdpGetTransactionId() int {
    return rand.Intn(255)
}


func UdpParseConnectionResponse(buf []byte, sent_transaction_id int32) (int64, error) {
    if len(buf) < 16 {
        return 0, fmt.Errorf("Wrong response length getting connection id: %d", len(buf))
    }

    action := int32(binary.BigEndian.Uint32(buf[:4])) // first 4 bytes is action
    res_transaction_id := int32(binary.BigEndian.Uint32(buf[4:8])) // next 4 bytes is transaction id

    if res_transaction_id != sent_transaction_id {
        return 0, fmt.Errorf("Transaction ID doesn't match in connection response! Expected %d, got %d", sent_transaction_id, res_transaction_id)
    }

    if action == 0 {
        connection_id := int64(binary.BigEndian.Uint64(buf[8:])) // unpack 8 bytes from byte 8, should be the connection_id
        return connection_id, nil
    } else if action == 3 {
        error_message := string(buf[8:])
        return 0, fmt.Errorf("Error while trying to get a connection response: %s", error_message)
    }

    return 0, fmt.Errorf("Unknown action type: %d", action)
}



func CreateUdpAnnounceRequest(connection_id int64, hashes []string) ([]byte, int) {
    action := uint32(0x1) // to announce
    transaction_id := uint32(UdpGetTransactionId())
    buf := make([]byte, 0)
    buf = append(buf, make([]byte, 8)...)
    binary.BigEndian.PutUint64(buf, uint64(connection_id))
    buf = append(buf, make([]byte, 4)...)
    binary.BigEndian.PutUint32(buf[len(buf)-4:], action)
    buf = append(buf, make([]byte, 4)...)
    binary.BigEndian.PutUint32(buf[len(buf)-4:], transaction_id)
    hex_repr, _ := hex.DecodeString(hashes[0])
    buf = append(buf, make([]byte, 20)...)
    copy(buf[len(buf)-20:], hex_repr)
    peer_id := "-MY"
    for i := 0; i < 4; i++ {
        peer_id += string(rand.Intn(10) + '0')
    }
    peer_id += "-"
    for i := 0; i < 12; i++ {
        peer_id += string(rand.Intn(10) + '0')
    }
    buf = append(buf, peer_id...)
    down := uint64(0x00)
    up := uint64(0x00)
    left := uint64(0x00)
    buf = append(buf, make([]byte, 8)...)
    binary.BigEndian.PutUint64(buf[len(buf)-8:], down)
    buf = append(buf, make([]byte, 8)...)
    binary.BigEndian.PutUint64(buf[len(buf)-8:], left)
    buf = append(buf, make([]byte, 8)...)
    binary.BigEndian.PutUint64(buf[len(buf)-8:], up)
    buf = append(buf, make([]byte, 4)...)
    binary.BigEndian.PutUint32(buf[len(buf)-4:], 0x2) // event 2 denotes start of downloading
    buf = append(buf, make([]byte, 4)...)
    binary.BigEndian.PutUint32(buf[len(buf)-4:], 0x0) // IP address set to 0. Response received to the sender of this packet
    key := uint32(UdpGetTransactionId()) // Unique key randomized by client
    buf = append(buf, make([]byte, 4)...)
    binary.BigEndian.PutUint32(buf[len(buf)-4:], key)
    buf = append(buf, make([]byte, 4)...)
    binary.BigEndian.PutUint32(buf[len(buf)-4:], 0xffffffff) // Number of peers required. Set to -1 for default
    buf = append(buf, make([]byte, 4)...)
    binary.BigEndian.PutUint32(buf[len(buf)-4:], 6882)
    return buf, int(transaction_id)
}


func udpParseAnnounceResponse(buf []byte, sentTransactionID int32) (map[string]int32, []map[string]interface{}, error) {
    if len(buf) < 20 {
        return nil, nil, fmt.Errorf("Wrong response length while announcing: %d", len(buf))
    }
    action := int32(binary.BigEndian.Uint32(buf[0:4]))
    resTransactionID := int32(binary.BigEndian.Uint32(buf[4:8]))
    if resTransactionID != sentTransactionID {
        return nil, nil, fmt.Errorf("Transaction ID doesnt match in announce response! Expected %d, got %d", sentTransactionID, resTransactionID)
    }
    if action == 0x1 {
        ret := make(map[string]int32)
        offset := 8
        ret["interval"] = int32(binary.BigEndian.Uint32(buf[offset : offset+4]))
        offset += 4
        ret["leeches"] = int32(binary.BigEndian.Uint32(buf[offset : offset+4]))
        offset += 4
        ret["seeds"] = int32(binary.BigEndian.Uint32(buf[offset : offset+4]))
        offset += 4
        var peers []map[string]interface{}
        for offset < len(buf)-4 {
            peer := make(map[string]interface{})
            peerIP := binary.BigEndian.Uint32(buf[offset : offset+4])
            peer["IP"] = net.IPv4(byte(peerIP>>24), byte(peerIP>>16), byte(peerIP>>8), byte(peerIP)).String()
            offset += 4
            if offset >= len(buf) {
                return nil, nil, fmt.Errorf("Error while reading peer port")
            }
            peer["port"] = binary.BigEndian.Uint16(buf[offset : offset+2])
            offset += 2
            peers = append(peers, peer)
        }
        return ret, peers, nil
    } else {
        errorStr := string(buf[8:])
        return nil, nil, fmt.Errorf("Error while announcing: %s", errorStr)
    }
}
