package common

import (
	"fmt"
	"os"
	"path/filepath"

	"torrent-dsp/model"
)

func CreateFile(torrent *model.Torrent) (*os.File, error) {
	fmt.Println("Creating folder structure for", torrent.Info.Name)

	// create the folder structure
	if torrent.Info.Files != nil {
		// multiple files
		for idx, filePath := range torrent.Info.Files {
			// Get the absolute path of the file
			absPath, err := filepath.Abs("downloads" + "/" + torrent.Info.Name + "/" + filePath.Path[0])
			if err != nil {
				return nil, err
			}
		
			// Get the directory path of the file
			dirPath := filepath.Dir(absPath)
		
			// Create the directory path if it doesn't exist
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return nil, err
			}
		
			// Create the file
			fmt.Println("Creating file", idx)
			file, err := os.Create(absPath)
			if err != nil {
				return nil, err
			}
			defer file.Close()
		}

		outFile, err := os.Create("downloads" + "/" + torrent.Info.Name + ".temp")
		if err != nil {
			return nil, err
		}
		return outFile, err
	} else {
		// single file
		outFile, err := os.Create("downloads" + "/" + torrent.Info.Name)
		if err != nil {
			return nil, err
		}
		return outFile, err
	}
}