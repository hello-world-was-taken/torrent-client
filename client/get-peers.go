package client

import "torrent-dsp/model"

import (
	"net"

)

func GetPeersFromTracker(torrent *model.Torrent, conn net.Conn) ([]*model.Peer, error) {
	// create a request to send to the tracker
	request, err := CreateTrackerRequest(torrent)
	if err != nil {
		return nil, err
	}

	// send the request to the tracker
	err = SendTrackerRequest(request, conn)
	if err != nil {
		return nil, err
	}

	// get the response from the tracker
	// response, err := GetTrackerResponse(conn)
	// if err != nil {
	// 	return nil, err
	// }

	// // parse the response and get a list of peers
	// peers, err := ParseTrackerResponse(response)
	// if err != nil {
	// 	return nil, err
	// }

	// return peers, nil
	return nil, nil
}

func CreateTrackerRequest(torrent *model.Torrent) (*model.TrackerRequest, error) {
	// create a new tracker request
	request := model.TrackerRequest{
		Info_hash:  torrent.InfoHash,
		Peer_id:    "-TR2770-6x6x6x6x6x6x",
		Port:       6889,
		Uploaded:   0,
		Downloaded: 0,
		Left:       torrent.Info.Length,
		Event:      "started",
	}

	return &request, nil
}


func SendTrackerRequest(request *model.TrackerRequest, conn net.Conn) error {
	// send the request to the tracker
	_, err := conn.Write([]byte(request.Encode()))
	if err != nil {
		return err
	}

	return nil
}

