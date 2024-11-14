package domain

import (
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

// ExtendedAccountRepository extends the imported AccountRepository by adding
// transactional methods, while keeping all original methods from AccountRepository.
type ExtendedAccountRepository interface {
	// Embed the original AccountRepository interface
	domain.AccountRepository

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}
