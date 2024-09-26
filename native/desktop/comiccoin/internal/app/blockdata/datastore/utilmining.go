package datastore

import (
	"context"
	"math/big"
)

func (b *Block) PerformPOW(ctx context.Context, difficulty uint16) error {

	// Choose a random starting point for the nonce. After this, the nonce
	// will be incremented by 1 until a solution is found by us or another node.
	nBig := big.NewInt(0)
	b.Header.Nonce = nBig.Uint64()

	for {
		// Did we timeout trying to solve the problem.
		if ctx.Err() != nil {
			return ctx.Err()
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
	return nil
}
