package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type ProofOfWorkUseCase struct {
	config *config.Config
	logger *slog.Logger
}

func NewProofOfWorkUseCase(config *config.Config, logger *slog.Logger) *ProofOfWorkUseCase {
	return &ProofOfWorkUseCase{config, logger}
}

func (uc *ProofOfWorkUseCase) Execute(ctx context.Context, b *domain.Block, difficulty uint16) (uint64, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if b == nil {
		e["block"] = "missing value"
	} else {
		if b.Header == nil {
			e["header"] = "missing value"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed executing proof of work",
			slog.Any("error", e))
		return 0, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	// Choose zero starting point for the nonce. After this, the nonce
	// will be incremented by 1 until a solution is found by us or another node.
	nBig := big.NewInt(0)
	b.Header.Nonce = nBig.Uint64()

	for {
		// Did we timeout trying to solve the problem.
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}

		// Hash the block and check if we have solved the puzzle.
		hash := b.Hash()
		if !isHashSolved(b.Header.Difficulty, hash) {
			b.Header.Nonce++
			continue
		} else {
			break
		}
	}

	return b.Header.Nonce, nil
}

// isHashSolved checks the hash to make sure it complies with
// the POW rules. We need to match a difficulty number of 0's.
func isHashSolved(difficulty uint16, hash string) bool {
	const match = "0x00000000000000000"

	if len(hash) != 66 {
		return false
	}

	difficulty += 2
	return hash[:difficulty] == match[:difficulty]
}
