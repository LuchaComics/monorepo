package static

import (
	_ "embed"
	"encoding/json"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

//go:embed genesis.json
var GenesisJSON string

// GetGenesisBlockData function returns the Genesis block data which was
// embedded in this application via the `genesis.json` file. The purpose of
// this function is to allow any third-party library to access directly the
// read-only contents of the Genesis Block.
func GetGenesisBlockData() (*domain.GenesisBlockData, error) {
	var genesis *domain.GenesisBlockData
	err := json.Unmarshal([]byte(GenesisJSON), &genesis)
	return genesis, err
}
