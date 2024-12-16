package main

import (
	_ "time/tzdata" // Important b/c some servers don't allow access to timezone file.

	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS to match Linux container CPU quota.

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/cmd"
)

func main() {
	cmd.Execute()
}
