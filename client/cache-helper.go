package client

import (
	"encoding/json"
	"fmt"
	"os"

	"torrent-dsp/model"
)


func SaveCache(filename string, cache *model.PiecesCache) error {
	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("error different from non existent")
		return err
	}

	// if the error is that the file does not exist, then create a new cache
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error opening cache while non existent")
			return err
		}
		file.Close()
	}
	file, err := os.OpenFile(filename, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Error opening cache")
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cache)
	if err != nil {
		fmt.Println("Error encoding cache", err)
		return err
	}

	return nil
}


// LoadCache loads the cache from a file
func LoadCache(filename string) (*model.PiecesCache, error) {
	// load the cache from a file
	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("error different from nil and err != from not exist")
		return nil, err
	}

	// if the error is that the file does not exist, then create a new cache
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error while creating file")
			return nil, err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.Encode(&model.PiecesCache{Pieces: map[int]bool{}})
		return &model.PiecesCache{Pieces: map[int]bool{}}, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}

	decoder := json.NewDecoder(file)
	var cache model.PiecesCache
	err = decoder.Decode(&cache)
	if err != nil {
		fmt.Println("error decoding file")
		return nil, err
	}

	return &cache, nil
}
