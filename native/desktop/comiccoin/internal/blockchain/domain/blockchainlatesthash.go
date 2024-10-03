package domain

type BlockchainLastestHashRepository interface {
	Get() (string, error)
	Set(hash string) error
}
