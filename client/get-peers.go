package client

import (
	"fmt"
	"net/http"
	"net/url"
	"torrent-dsp/model"
	"torrent-dsp/utils"

	bencode "github.com/zeebo/bencode"
)

// Get peers from one of the trackers in the torrent file
func GetPeersFromTrackers(torrent *model.Torrent) ([]*model.Peer, error) {
	// create a request to send to the tracker
	httpTrackerURLs, err := buildTrackerRequestURLs(torrent)
	if err != nil {
		return nil, err
	}

	go getPeersFromTrackersHelper(httpTrackerURLs)
	if err != nil {
		return nil, err
	}

	// return peers, nil
	return nil, nil
}


// builds a tracker request url from the announcement list
func buildTrackerRequestURLs(torrent *model.Torrent) ([]string, error) {
	// query parameters to the tracker url
	requestParams := model.TrackerRequestParams{
		Info_hash:  torrent.InfoHash,
		Peer_id:    "-TR2940-k8hj0wgej6ch",
		Port:       6881,
		Uploaded:   0,
		Downloaded: 0,
		Left:       torrent.Info.Length,
		Event:      "started",
	}

	// list of urls
	URLs := []string{}

	// go through all the trackers in the announce list and create a request for each one
	// discard udp and only take http trackers
	for _, tracker := range torrent.AnnounceList {
		// parse the tracker url
		URL, err := url.Parse(tracker[0])
		if err != nil {
			return []string{}, err
		}

		if utils.IsHTTPTracker(URL.String()) {
			// add the query parameters to the tracker url
			URL.RawQuery = requestParams.Encode()
			URLs = append(URLs, URL.String())
		}
	}

	return URLs, nil
}


// go through each HTTP tracker requests until we get a response from one
func getPeersFromTrackersHelper(URLs []string) (model.TrackerResponse, error) {

	for _, URL := range URLs {
		response, err := getPeerFromURL(URL)
		if err == nil {
			return response, nil
		}

		fmt.Println("Error getting peers from HTTP tracker: ", err)
	}

	return model.TrackerResponse{}, nil
}


// try to get peers from a passed in tracker request url
func getPeerFromURL(URL string) (model.TrackerResponse, error) {
	// send the request to the tracker
	response, err := http.Get(URL)
	if err != nil {
		return model.TrackerResponse{}, err
	}
	
	// close the response body
	defer response.Body.Close()

	// decode the response
	trackerResponse := model.TrackerResponse{}
	err = bencode.NewDecoder(response.Body).Decode(&trackerResponse)
	if err != nil {
		return model.TrackerResponse{}, err
	}

	fmt.Println("Retrieved peers successfully")
	return trackerResponse, nil
}

