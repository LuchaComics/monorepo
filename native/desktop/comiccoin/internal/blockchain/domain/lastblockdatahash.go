package domain

type LastBlockDataHashRepository interface {
	Get() (string, error)
	Set(hash string) error
}
