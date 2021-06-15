package integrity

import "time"

type IntegrityFile struct {
	TargetPath            string
	ModificationTimeStamp time.Time
	Entries               []IntegrityEntry
}

func (file *IntegrityFile) EntriesMap() map[string]IntegrityEntry {
	var integrityMap = make(map[string]IntegrityEntry)
	for _, integrityEntry := range file.Entries {
		integrityMap[integrityEntry.Path] = integrityEntry
	}
	return integrityMap
}

type IntegrityEntry struct {
	Path       string
	Checksum   string
	ChangeDate time.Time
}
