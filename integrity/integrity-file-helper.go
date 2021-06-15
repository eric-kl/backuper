package integrity

import (
	"encoding/json"
	"errors"
	"os"
)

func WriteFile(file IntegrityFile, filePath string) error {
	res, err := json.MarshalIndent(file, "", "  ")

	if err != nil {
		return errors.New("Could not serialize integrity File!")
	}

	err = os.WriteFile(filePath, res, 0777)

	if err != nil {
		return errors.New("Could not write integrity File!")
	}

	return nil
}

func ReadFile(path string) (IntegrityFile, error) {
	res, err := os.ReadFile(path)

	if err != nil {
		return IntegrityFile{}, errors.New("Could not Read Integrity File")
	}

	var fileContent IntegrityFile

	err = json.Unmarshal(res, &fileContent)

	if err != nil {
		return IntegrityFile{}, errors.New("Could not deserialize Integrity File")
	}

	return fileContent, nil
}
