package constants

const (
	DefaultIdentityKeyID = "blockchain-node"
)

// Distributed publish-subscribe broker constants
const (
	PubSubMempoolTopicName = "mempool"
)

const (
	ChainIDMainNet = 1
)

const (
	ConsensusPoW = "PoW"
	ConsensusPoA = "PoA"
)

// Unified constant values to use for all ComicCoin repositories.
const (
	ComicCoinChainID                        = ChainIDMainNet
	ComicCoinTransPerBlock                  = 1
	ComicCoinDifficulty                     = 2
	ComicCoinConsensusPollingDelayInMinutes = 1
	ComicCoinConsensusProtocol              = ConsensusPoA
	ComicCoinPeerListenPort                 = 26644
	ComicCoinBootstrapPeers                 = "/ip4/127.0.0.1/tcp/26642/p2p/QmfYJfA7uRJNodrozCLAB2VWcp8dP8q289E8JwhJY5HSi7"
	ComicCoinIdentityKeyID                  = DefaultIdentityKeyID
	ComicCoinNFTAssetStoreAddress           = "http://127.0.0.1:8080" // TODO: Change to PROD NFT asset store when ready.
)
