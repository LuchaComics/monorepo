package datastore

import (
	"context"
	"io/ioutil"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func (impl *accountStorerImpl) Insert(ctx context.Context, b *Account) error {
	bBytes, err := b.Serialize()
	if err != nil {
		return err
	}
	if err := impl.dbClient.Setf(bBytes, "account-%v", b.Name); err != nil {
		return err
	}
	return nil
}

func (impl *accountStorerImpl) GetByName(ctx context.Context, name string) (*Account, error) {
	bBytes, err := impl.dbClient.Getf("account-%v", name)
	if err != nil {
		return nil, err
	}
	b, err := NewAccountFromDeserialize(bBytes)
	if err != nil {
		impl.logger.Error("failed to deserialize",
			slog.String("name", name),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (impl *accountStorerImpl) List(ctx context.Context) ([]*Account, error) {
	res := make([]*Account, 0)
	seekThenIterateKey := ""
	err := impl.dbClient.Iterate("account-", seekThenIterateKey, func(key, value []byte) error {
		account, err := NewAccountFromDeserialize(value)
		if err != nil {
			impl.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, account)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (impl *accountStorerImpl) DeleteByName(ctx context.Context, name string) error {
	err := impl.dbClient.Deletef("account-%v", name)
	if err != nil {
		return err
	}
	return nil
}

func (impl *accountStorerImpl) GetKeyByNameAndPassword(ctx context.Context, accountName string, accountWalletPassword string) (*keystore.Key, error) {
	account, err := impl.GetByName(ctx, accountName)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil
	}

	keyJson, err := ioutil.ReadFile(account.WalletFilepath)
	if err != nil {
		return nil, nil
	}

	key, err := keystore.DecryptKey(keyJson, accountWalletPassword)
	if err != nil {
		return nil, nil
	}
	return key, nil
}
