package integrity

type IntegrityDif struct {
	AddedFiles   []IntegrityEntry
	RemovedFiles []IntegrityEntry
	ChangedFiles []IntegrityEntry
}
