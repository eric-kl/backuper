package integrity

import (
	"os"
	"path/filepath"
	"time"
)

func UpdateIntegrityFile(dif *IntegrityDif, integrityFile *IntegrityFile) IntegrityFile {
	//Map file entries to map for fast and easy manipulation
	//var integrityMap = make(map[string]IntegrityEntry)
	//for _, integrityEntry := range integrityFile.Entries {
	//	integrityMap[integrityEntry.Path] = integrityEntry
	//}
	integrityMap := integrityFile.EntriesMap()

	//remove lost files
	for _, entry := range dif.RemovedFiles {
		delete(integrityMap, entry.Path)
	}

	//Add added files
	for _, entry := range dif.AddedFiles {
		integrityMap[entry.Path] = entry
	}

	//Override updated files
	for _, entry := range dif.ChangedFiles {
		integrityMap[entry.Path] = entry
	}

	integrityFile.ModificationTimeStamp = dif.CreationTime
	integrityFile.Entries = []IntegrityEntry{}
	for _, entry := range integrityMap {
		integrityFile.Entries = append(integrityFile.Entries, entry)
	}

	return *integrityFile
}

func GetIntegrityDif(pathToTarget string, integrityFile *IntegrityFile) (IntegrityDif, error) {
	var dif = IntegrityDif{CreationTime: time.Now()}

	//var integrityMap = make(map[string]IntegrityEntry)
	//for _, integrityEntry := range integrityFile.Entries {
	//	integrityMap[integrityEntry.Path] = integrityEntry
	//}
	integrityMap := integrityFile.EntriesMap()

	err := filepath.Walk(pathToTarget, func(pathToFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		//Turn path into a relative path
		relativePath := pathToFile[len(pathToTarget):]
		osAgnosticPath := filepath.ToSlash(relativePath)

		entry, exists := integrityMap[osAgnosticPath]
		if !exists {
			//Add to added list if not exists
			hash, err := hash_file_md5(pathToFile)

			if err != nil {
				return err
			}

			dif.AddedFiles = append(dif.AddedFiles, IntegrityEntry{Path: osAgnosticPath, Checksum: hash, ChangeDate: info.ModTime()})
			return nil
		}

		//Remove entry from integrityMap it has been found in the new filepath
		delete(integrityMap, entry.Path)

		//First check the date
		if entry.ChangeDate == info.ModTime() {
			return nil //Same Date means same file (for soft integrity)
		}

		hash, err := hash_file_md5(pathToFile)

		if err != nil {
			return err
		}

		if hash == entry.Checksum {
			return nil //Its actually still the same file, no write neccessary
		}

		dif.ChangedFiles = append(dif.ChangedFiles, IntegrityEntry{Path: osAgnosticPath, Checksum: hash, ChangeDate: info.ModTime()})
		return nil
	})

	if err != nil {
		return IntegrityDif{}, err
	}

	for _, entry := range integrityMap {
		dif.RemovedFiles = append(dif.RemovedFiles, entry)
	}

	return dif, nil
}

func DifBetweenIntegrityFiles(fileA *IntegrityFile, fileB *IntegrityFile) IntegrityDif {
	//var integrityMap = make(map[string]IntegrityEntry)
	//for _, integrityEntry := range fileA.Entries {
	//	integrityMap[integrityEntry.Path] = integrityEntry
	//}

	integrityMap := fileA.EntriesMap()

	var dif = IntegrityDif{CreationTime: time.Now()}

	for _, entryB := range fileB.Entries {
		entryA, exists := integrityMap[entryB.Path]
		if !exists {
			dif.AddedFiles = append(dif.AddedFiles, entryB)
			continue
		}

		if entryA.ChangeDate != entryB.ChangeDate && entryA.Checksum != entryB.Checksum {
			dif.ChangedFiles = append(dif.ChangedFiles, entryB)
		}
		delete(integrityMap, entryA.Path)
	}

	for _, entryA := range integrityMap {
		dif.RemovedFiles = append(dif.RemovedFiles, entryA)
	}

	return dif
}
