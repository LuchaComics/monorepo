package domain

type LastBlockDataHash string

type LastBlockDataHashRepository interface {
	Get() (LastBlockDataHash, error)
	Set(hash LastBlockDataHash) error
}
