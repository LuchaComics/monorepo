package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ProofOfWorkUseCase struct {
	config *config.Configuration
	logger *slog.Logger
}

func NewProofOfWorkUseCase(config *config.Configuration, logger *slog.Logger) *ProofOfWorkUseCase {
	return &ProofOfWorkUseCase{config, logger}
}

func (uc *ProofOfWorkUseCase) Execute(ctx context.Context, b *domain.Block, difficulty uint16) (*big.Int, error) {
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
		return big.NewInt(0), httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	// Choose zero starting point for the nonce. After this, the nonce
	// will be incremented by 1 until a solution is found by us or another node.
	nBig := big.NewInt(0)
	b.Header.NonceBytes = nBig.Bytes()

	for {
		// Did we timeout trying to solve the problem.
		if ctx.Err() != nil {
			return big.NewInt(0), ctx.Err()
		}

		// Hash the block and check if we have solved the puzzle.
		hash := b.Hash()
		if !isHashSolved(b.Header.Difficulty, hash) {
			nBig.Add(nBig, big.NewInt(1))
			b.Header.SetNonce(nBig)
			continue
		} else {
			break
		}
	}

	return b.Header.GetNonce(), nil
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
