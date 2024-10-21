package constants

const (
	PageIDExit = iota
	PageIDPickDataDirectoryView
	PageIDStartupView
	PageIDMainAppView
)

const (
	TabIDOverivew = iota
	TabIDSend
	TabIDReceive
	TabIDTransactions
	TabIDSettings
)

const (
	DefaultDataDirectoryPath = "./ComicCoin"
	DefaultScreenWidth       = 640
	DefaultScreenHeight      = 480
)

const (
	PreferenceKeyHasSetDataDirectory = "HasSetDataDirectory"
	PreferenceKeyDataDirectory       = "DataDirectory"
)
