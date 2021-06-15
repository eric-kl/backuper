package main

import (
	"backuper/config"
	"backuper/integrity"
	"flag"
	"fmt"
	"os"
)

func main() {
	var update bool
	flag.BoolVar(&update, "update", false, "updates the integrity file. Default false")

	flag.Parse()

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

	fmt.Println("Read Config file completed")
	integrityFile, err := readOrCreateIntegrityFile(config.IntegrityPath, config.Target)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if update {
		err = updateAndWriteIntegrityFile(config.Target, config.IntegrityPath, &integrityFile)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	//if no update just show dif
	dif, err := integrity.GetIntegrityDif(config.Target, &integrityFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	isFilesAdded := len(dif.AddedFiles) != 0
	isFilesRemoved := len(dif.RemovedFiles) != 0
	isFilesChanged := len(dif.ChangedFiles) != 0

	if !isFilesAdded && !isFilesChanged && !isFilesRemoved {
		fmt.Println("No Changes, Integrity is up to date")
	}

	if isFilesAdded {
		fmt.Println("Added Files:")
		for _, entry := range dif.AddedFiles {
			fmt.Println(entry.Path)
		}
	}

	if isFilesRemoved {
		fmt.Println("Removed Files:")
		for _, entry := range dif.RemovedFiles {
			fmt.Println(entry.Path)
		}
	}

	if isFilesChanged {
		fmt.Println("Changed Files:")
		for _, entry := range dif.ChangedFiles {
			fmt.Println(entry.Path)
		}
	}

}

func updateAndWriteIntegrityFile(pathToTarget string, pathToIntegrityFile string, integrityFile *integrity.IntegrityFile) error {
	newIntegrityFile, err := integrity.UpdateIntegrityFile(pathToTarget, integrityFile)

	if err != nil {
		return err
	}

	err = integrity.WriteFile(newIntegrityFile, pathToIntegrityFile)

	if err != nil {
		return err
	}

	fmt.Println("Wrote integrity file")
	return nil
}

func readOrCreateIntegrityFile(pathToIntegrityFile string, targetPath string) (integrity.IntegrityFile, error) {
	_, err := os.Stat(pathToIntegrityFile)
	if err != nil {
		fmt.Println("no integrity file found at " + pathToIntegrityFile + " will continue with an empty file")
		return integrity.IntegrityFile{}, nil
	}

	return integrity.ReadFile(pathToIntegrityFile)
}
