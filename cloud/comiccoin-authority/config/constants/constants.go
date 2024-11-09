package constants

type key int

const (
	SessionIsAuthorized key = iota
	SessionSkipAuthorization
	SessionID
	SessionIPAddress
	SessionProxies
	SessionUser
	SessionUserCompanyName
	SessionUserRole
	SessionUserID
	SessionUserUUID
	SessionUserTimezone
	SessionUserName
	SessionUserFirstName
	SessionUserLastName
	SessionUserStoreID
	SessionUserStoreName
	SessionUserStoreLevel
	SessionUserStoreTimezone
)

const (
	DefaultIdentityKeyID = "blockchain-node"
)

// Distributed publish-subscribe broker constants
const (
	PubSubMempoolTopicName = "mempool"
)

const (
	ChainIDMainNet = 1
	ChainIDTestNet = 2
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
)
