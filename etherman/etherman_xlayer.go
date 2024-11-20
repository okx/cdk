package etherman

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/0xPolygon/cdk/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// LoadAuthFromKeyStoreXLayer loads an authorization from a key store file
func (etherMan *Client) LoadAuthFromKeyStoreXLayer(path, password string) (
	*bind.TransactOpts,
	*ecdsa.PrivateKey,
	error,
) {
	auth, pk, err := newAuthFromKeystoreXLayer(path, password, etherMan.l1Cfg.L1ChainID)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("loaded authorization for address: %v", auth.From.String())
	etherMan.auth[auth.From] = auth
	return &auth, pk, nil
}

// newAuthFromKeystoreXLayer an authorization instance from a keystore file
func newAuthFromKeystoreXLayer(path, password string, chainID uint64) (bind.TransactOpts, *ecdsa.PrivateKey, error) {
	log.Infof("reading key from: %v", path)
	key, err := newKeyFromKeystore(path, password)
	if err != nil {
		return bind.TransactOpts{}, nil, err
	}
	if key == nil {
		return bind.TransactOpts{}, nil, nil
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		return bind.TransactOpts{}, nil, err
	}
	return *auth, key.PrivateKey, nil
}
