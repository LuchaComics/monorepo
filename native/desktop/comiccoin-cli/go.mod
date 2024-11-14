module github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli

go 1.23.0

require (
	github.com/LuchaComics/monorepo/cloud/comiccoin-authority v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/LuchaComics/monorepo/cloud/comiccoin-authority => ../../../cloud/comiccoin-authority
