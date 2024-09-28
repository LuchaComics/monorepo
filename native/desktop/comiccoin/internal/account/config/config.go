package config

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	// DirPath variable is the path to the directory where all the files for
	// this appliction to
	// save to.
	DirPath string

	// HTTPAddress variable is the address and port that the HTTP JSON API
	// server will listen on for this application. Do not expose to public!
	HTTPAddress string

	// RPCAddress variable is the address and port that the TCP RCP
	// server will listen on for this application. Do not expose to public!
	RPCAddress string
}

type DBConfig struct {
	// Location of were to save the database files.
	DataDir string
}
