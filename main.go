package main

import (
	"backuper/config"
	"backuper/integrity"
	"fmt"
	"os"
)

func main() {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("No working directory")
		return
	}
	fmt.Println("WorkingDirectory found as ", currentWorkingDirectory)

	config, err := config.ReadConfig(currentWorkingDirectory)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("start creating integrity for " + config.Target)
	file, err := integrity.CreateIntegrityFileForPath(config.Target)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = integrity.WriteFile(file, config)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Wrote integrity file")
}
