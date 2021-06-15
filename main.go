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

	dif, err := integrity.GetDifAndPrint(&integrityFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	if update {
		err = updateAndWriteIntegrityFile(&dif, config.IntegrityPath, &integrityFile)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func updateAndWriteIntegrityFile(dif *integrity.IntegrityDif, pathToIntegrityFile string, integrityFile *integrity.IntegrityFile) error {
	newIntegrityFile := integrity.UpdateIntegrityFile(dif, integrityFile)

	err := integrity.WriteFile(newIntegrityFile, pathToIntegrityFile)

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
		return integrity.CreateFile(targetPath), nil
	}

	return integrity.ReadFile(pathToIntegrityFile)
}
