package integrity

type IntegrityFile struct {
	Entries []IntegrityEntry
}

type IntegrityEntry struct {
	Path     string
	Checksum string
}
