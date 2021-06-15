package config

type ConfigFile struct {
	Target        string
	IntegrityPath string
	//Soft Integrity true means that the same mod date is viewed as the same file for integrity,
	//false means always use hash to compare files
	SoftIntegrity bool
}
