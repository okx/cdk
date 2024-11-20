package etherman

import (
	"errors"
	"fmt"
	"math/big"

	polygonzkevm "github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonvalidiumetrog"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BuildMockSequenceBatchesTxData builds a []bytes to be sent to the PoE SC method SequenceBatches.
func (etherMan *Client) BuildMockSequenceBatchesTxData(sender common.Address,
	validiumBatchData []polygonzkevm.PolygonValidiumEtrogValidiumBatchData,
	maxSequenceTimestamp uint64,
	l2Coinbase common.Address,
	dataAvailabilityMessage []byte,
	l1InfoTreeLeafCount uint32,
	expectedInputHash string,
) (to *common.Address, data []byte, err error) {
	opts, err := etherMan.getAuthByAddress(sender)
	if errors.Is(err, ErrNotFound) {
		return nil, nil, fmt.Errorf("failed to build sequence batches, err: %w", ErrPrivateKeyNotFound)
	}
	opts.NoSend = true
	// force nonce, gas limit and gas price to avoid querying it from the chain
	opts.Nonce = big.NewInt(1)
	opts.GasLimit = uint64(1)
	opts.GasPrice = big.NewInt(1)

	tx, err := etherMan.sequenceMockBatches(
		opts,
		validiumBatchData,
		maxSequenceTimestamp,
		l2Coinbase,
		dataAvailabilityMessage,
		l1InfoTreeLeafCount,
		expectedInputHash,
	)
	if err != nil {
		return nil, nil, err
	}

	return tx.To(), tx.Data(), nil
}

func (etherMan *Client) sequenceMockBatches(opts bind.TransactOpts,
	validiumBatchData []polygonzkevm.PolygonValidiumEtrogValidiumBatchData,
	maxSequenceTimestamp uint64,
	l2Coinbase common.Address,
	dataAvailabilityMessage []byte,
	l1InfoTreeLeafCount uint32,
	expectedInputHash string,
) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error

	var finalInputHash = [32]byte{}
	// Request will a hex string beginning with "0x...", so strip first 2 chars.
	for i, bb := range common.Hex2Bytes(expectedInputHash[2:]) {
		finalInputHash[i] = bb
	}

	tx, err = etherMan.Contracts.Banana.Rollup.SequenceBatchesValidium(
		&opts,
		validiumBatchData,
		l1InfoTreeLeafCount,
		maxSequenceTimestamp,
		finalInputHash,
		l2Coinbase,
		dataAvailabilityMessage,
	)

	if err != nil {
		if parsedErr, ok := TryParseError(err); ok {
			err = parsedErr
		}
		err = fmt.Errorf(
			"error sequencing batches: %w, dataAvailabilityMessage: %s",
			err, common.Bytes2Hex(dataAvailabilityMessage),
		)
	}

	return tx, err
}
