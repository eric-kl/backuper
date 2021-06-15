package integrity

import (
	"backuper/config"
	"encoding/json"
	"errors"
	"os"
)

func WriteFile(file IntegrityFile, config config.ConfigFile) error {
	res, err := json.MarshalIndent(file, "", "  ")

	if err != nil {
		return errors.New("Could not serialize integrity File!")
	}

	err = os.WriteFile(config.IntegrityPath, res, 0777)

	if err != nil {
		return errors.New("Could not write integrity File!")
	}

	return nil
}
