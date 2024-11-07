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
	ComicCoinPeerListenPort                 = 26644
	ComicCoinBootstrapPeers                 = "/ip4/127.0.0.1/tcp/26642/p2p/QmZi2NPQx41oxfXdpVNGFtWp4rN1RJoYVyvAinKM9MHLh1" // Example `/ip4/127.0.0.1/tcp/26642/p2p/QmXYZ`.
	ComicCoinIdentityKeyID                  = constants.DefaultIdentityKeyID
	ComicCoinIPFSPublicGatewayDomain        = "http://127.0.0.1:8080" // TODO: Change to PROD when ready
)
