package txbuilder

import (
	"errors"
	"fmt"

	"github.com/0xPolygon/cdk/etherman"
	rpctypes "github.com/0xPolygon/cdk/rpc/types"
)

// RPCInterface represents the RPC interface
type RPCInterface interface {
	GetBatch(batchNumber uint64) (*rpctypes.RPCBatch, error)
}

func (t *TxBuilderBananaValidium) checkMaxTimestamp(sequence etherman.SequenceBanana) error {
	if t.l2RpcClient == nil {
		return nil
	}
	var maxBatchNumber uint64
	for _, batch := range sequence.Batches {
		maxBatchNumber = max(maxBatchNumber, batch.BatchNumber)
	}
	rpcBatch, err := t.l2RpcClient.GetBatch(maxBatchNumber)
	if err != nil {
		t.logger.Error("check rpc max timestamp error: ", err)
		return err
	}
	if rpcBatch.LastL2BLockTimestamp() != sequence.MaxSequenceTimestamp {
		t.logger.Error("max timestamp mismatch: ", rpcBatch.LastL2BLockTimestamp(), sequence.MaxSequenceTimestamp)
		return errors.New(fmt.Sprintf("max timestamp mismatch: %v, %v",
			rpcBatch.LastL2BLockTimestamp(), sequence.MaxSequenceTimestamp))
	}
	t.logger.Infof("max timestamp check passed:%v,%v", maxBatchNumber, sequence.MaxSequenceTimestamp)
	return nil
}
