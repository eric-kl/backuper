package integrity

import "time"

type IntegrityDif struct {
	CreationTime time.Time
	AddedFiles   []IntegrityEntry
	RemovedFiles []IntegrityEntry
	ChangedFiles []IntegrityEntry
}
