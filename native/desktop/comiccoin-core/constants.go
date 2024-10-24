package main

import (
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
)

const (
	ComicCoinChainID                        = constants.ChainIDMainNet
	ComicCoinTransPerBlock                  = 1
	ComicCoinDifficulty                     = 2
	ComicCoinConsensusPollingDelayInMinutes = 1
	ComicCoinConsensusProtocol              = constants.ConsensusPoA
	ComicCoinPeerListenPort                 = 26642
	ComicCoinBootstrapPeers                 = "/ip4/127.0.0.1/tcp/26642/p2p/QmeZtjicLB8wWJ9L1gPuD5wrjmKQwJU6gHGHoaSofkXtUx" // Example `/ip4/127.0.0.1/tcp/26642/p2p/QmXYZ`.
	ComicCoinIdentityKeyID                  = constants.DefaultIdentityKeyID
)
