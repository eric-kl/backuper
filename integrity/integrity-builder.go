package integrity

import (
	"os"
	"path/filepath"
)

func CreateIntegrityFileForPath(path string) (IntegrityFile, error) {
	return UpdateIntegrityFile(path, &IntegrityFile{})
}

func UpdateIntegrityFile(pathToTarget string, integrityFile *IntegrityFile) (IntegrityFile, error) {
	dif, err := GetIntegrityDif(pathToTarget, integrityFile)

	if err != nil {
		return IntegrityFile{}, err
	}

	//Map file entries to map for fast and easy manipulation
	var integrityMap = make(map[string]IntegrityEntry)
	for _, integrityEntry := range integrityFile.Entries {
		integrityMap[integrityEntry.Path] = integrityEntry
	}

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

	integrityFile.Entries = []IntegrityEntry{}
	for _, entry := range integrityMap {
		integrityFile.Entries = append(integrityFile.Entries, entry)
	}

	return *integrityFile, nil
}

func GetIntegrityDif(pathToTarget string, integrityFile *IntegrityFile) (IntegrityDif, error) {
	var dif IntegrityDif

	var integrityMap = make(map[string]IntegrityEntry)
	for _, integrityEntry := range integrityFile.Entries {
		integrityMap[integrityEntry.Path] = integrityEntry
	}

	err := filepath.Walk(pathToTarget, func(pathToFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		entry, exists := integrityMap[pathToFile]
		if !exists {
			//Add to added list if not exists
			hash, err := hash_file_md5(pathToFile)

			if err != nil {
				return err
			}

			dif.AddedFiles = append(dif.AddedFiles, IntegrityEntry{Path: pathToFile, Checksum: hash, ChangeDate: info.ModTime()})
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

		dif.ChangedFiles = append(dif.ChangedFiles, IntegrityEntry{Path: pathToFile, Checksum: hash, ChangeDate: info.ModTime()})
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
