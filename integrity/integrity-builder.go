package integrity

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func CreateIntegrityFileForPath(path string) (IntegrityFile, error) {
	var entries []IntegrityEntry

	err := filepath.Walk(path, func(pathToFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		hash, err := hash_file_md5(pathToFile)

		if err != nil {
			return err
		}

		entries = append(entries, IntegrityEntry{Path: pathToFile, Checksum: hash})

		return nil
	})

	if err != nil {
		return IntegrityFile{}, err
	}

	return IntegrityFile{Entries: entries}, nil
}

func hash_file_md5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}
