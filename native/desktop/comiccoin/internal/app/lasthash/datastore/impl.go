package datastore

import "context"

func (impl *lastHashStorerImpl) Get(ctx context.Context) (string, error) {
	return "", nil
}

func (impl *lastHashStorerImpl) Set(ctx context.Context, hash string) error {
	return nil
}
