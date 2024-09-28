package domain

// BlockTransaction represents the transaction as it's recorded inside a block. This
// includes a timestamp and gas fees.
type BlockTransaction struct {
	SignedTransaction
	TimeStamp uint64 `json:"timestamp"` // Ethereum: The time the transaction was received.
	GasPrice  uint64 `json:"gas_price"` // Ethereum: The price of one unit of gas to be paid for fees.
	GasUnits  uint64 `json:"gas_units"` // Ethereum: The number of units of gas used for this transaction.
}
