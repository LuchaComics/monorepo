package domain

import (
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

// ExtendedWalletRepository extends the imported WalletRepository by adding
// transactional methods, while keeping all original methods from WalletRepository.
type ExtendedWalletRepository interface {
	// Embed the original WalletRepository interface
	domain.WalletRepository

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}
