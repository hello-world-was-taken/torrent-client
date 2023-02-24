package common

import (
	"fmt"
	"os"
	"path/filepath"

	"torrent-dsp/model"
)

func CreateFile(torrent *model.Torrent) (*os.File, string, error) {
	fmt.Println("Creating folder structure for", torrent.Info.Name)

	// create the folder structure
	if torrent.Info.Files != nil {
		// multiple files
		for idx, filePath := range torrent.Info.Files {
			// Get the absolute path of the file
			absPath, err := filepath.Abs("downloads" + "/" + torrent.Info.Name + "/" + filePath.Path[0])
			if err != nil {
				return nil, "", err
			}

			// Get the directory path of the file
			dirPath := filepath.Dir(absPath)

			// Create the directory path if it doesn't exist
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return nil, "", err
			}

			// Create the file
			fmt.Println("Creating file", idx)
			file, err := os.Create(absPath)
			if err != nil {
				return nil, "", err
			}
			defer file.Close()
		}

		outFile, err := CreateOrOpenFile("downloads" + "/" + torrent.Info.Name + ".temp")
		if err != nil {
			return nil, "", err
		}
		return outFile, "", err
	} else {
		// single file
		outFile, err := CreateOrOpenFile("downloads" + "/" + torrent.Info.Name)
		if err != nil {
			return nil, "", err
		}
		return outFile, "", err
	}
}

func CreateOrOpenFile(filename string) (*os.File, error) {
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
			fmt.Println("Error while creating file", filename)
			return nil, err
		}
		return file, nil
	}

	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}

	return file, nil
}
