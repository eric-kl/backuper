package integrity

import "time"

type IntegrityFile struct {
	Entries []IntegrityEntry
}

type IntegrityEntry struct {
	Path       string
	Checksum   string
	ChangeDate time.Time
}
