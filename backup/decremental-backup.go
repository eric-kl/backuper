package backup

import (
	"backuper/integrity"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

func CreateDerementalBackup(integrityFileBackup *integrity.IntegrityFile, integrityFileSource *integrity.IntegrityFile) error {
	dif := integrity.DifBetweenIntegrityFiles(integrityFileBackup, integrityFileSource)
	pathToChangesets := ""
	backupTarget := ""
	sourceTarget := ""

	for _, entry := range dif.RemovedFiles {
		//move entry to changeset
		err := moveFile(backupTarget, pathToChangesets, entry.Path)
		if err != nil {
			return errors.New("Could not move file")
		}
	}

	for _, entry := range dif.ChangedFiles {
		//move entry to changeset
		err := moveFile(backupTarget, pathToChangesets, entry.Path)
		if err != nil {
			return errors.New("Could not move file")
		}

		//move modified file to backup
		err = moveFile(sourceTarget, backupTarget, entry.Path)
		if err != nil {
			return errors.New("Could not move file")
		}
	}

	for _, entry := range dif.AddedFiles {
		//copy entry to backup
		oldPath := getFilesystemPath(sourceTarget, entry.Path)
		newPath := getFilesystemPath(backupTarget, entry.Path)
		in, err := os.Open(oldPath)

		if err != nil {
			return errors.New("Unable to open file")
		}
		defer in.Close()

		out, err := os.Create(newPath)
		if err != nil {
			return errors.New("Unable to Create file")
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return errors.New("Could not copy file")
		}
	}

	//TODO: Write Log to changeset
	return nil
}

func moveFile(sourcePath string, destPath string, relativeEntryPath string) error {
	oldPath := getFilesystemPath(sourcePath, relativeEntryPath)
	newPath := getFilesystemPath(destPath, relativeEntryPath)
	return os.Rename(oldPath, newPath)
}

func getFilesystemPath(directoryPath string, relativeEntryPath string) string {
	fullPath := path.Join(directoryPath, relativeEntryPath)
	return filepath.FromSlash(fullPath)
}
