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
	ComicCoinBootstrapPeers                 = "/ip4/127.0.0.1/tcp/26642/p2p/QmZi2NPQx41oxfXdpVNGFtWp4rN1RJoYVyvAinKM9MHLh1" // Example `/ip4/127.0.0.1/tcp/26642/p2p/QmXYZ`.
	ComicCoinIdentityKeyID                  = DefaultIdentityKeyID
	ComicCoinIPFSPublicGatewayDomain        = "http://127.0.0.1:8080" // TODO: Change to PROD when ready
)
