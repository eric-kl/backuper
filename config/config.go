package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func ReadConfig(path string) (ConfigFile, error) {
	configFileWithPath := path + "\\config.json"
	_, err := os.Stat(configFileWithPath)
	if err != nil {
		return createDefaultConfig(path)
	}

	content, err := ioutil.ReadFile(configFileWithPath)
	if err != nil {
		return ConfigFile{}, errors.New("Error when opening file: " + err.Error())
	}

	var fileContent ConfigFile

	err = json.Unmarshal(content, &fileContent)

	if err != nil {
		return ConfigFile{}, errors.New("Could not deserialize file: " + err.Error())
	}

	return fileContent, nil
}

func createDefaultConfig(path string) (ConfigFile, error) {
	configFileWithPath := path + "\\config.json"
	config := ConfigFile{Target: path, IntegrityPath: path + "\\integrity.json"}
	res, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		return ConfigFile{}, errors.New("Could not create config file!")
	}

	err = os.WriteFile(configFileWithPath, res, 0777)

	if err != nil {
		return ConfigFile{}, errors.New("Could not create config file!")
	}

	return config, errors.New("Config did not exist, created new Config, please restart after filling it out")
}
