package integrity

import (
	"fmt"
	"path/filepath"
	"time"
)

func GetDifAndPrint(integrityFile *IntegrityFile) (IntegrityDif, error) {
	targetPath := filepath.FromSlash(integrityFile.TargetPath)
	fmt.Println("Try to Get Dif for " + targetPath)
	startTime := time.Now()
	dif, err := GetIntegrityDif(targetPath, integrityFile)
	endTime := time.Now()

	if err != nil {
		//fmt.Println(err)
		return dif, err
	}

	totalChanges := len(dif.AddedFiles) + len(dif.RemovedFiles) + len(dif.ChangedFiles)

	if totalChanges == 0 {
		fmt.Println("No Changes, Integrity is up to date")
	}

	printFiles("Added Files: ", dif.AddedFiles)
	printFiles("Removed Files: ", dif.RemovedFiles)
	printFiles("Changed Files: ", dif.ChangedFiles)

	fmt.Println("Total Changes: ", totalChanges)
	fmt.Println("Found in: ", endTime.Sub(startTime).Round(time.Second))

	return dif, nil
}

func printFiles(preLine string, entries []IntegrityEntry) {
	if len(entries) == 0 {
		return
	}

	fmt.Println(preLine, len(entries))
	for _, entry := range entries {
		fmt.Println(entry.Path)
	}
}
