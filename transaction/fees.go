package transaction

import (
	"slices"

	"github.com/pkg/errors"
)

type ChangeDistribution int

var (
	ChangeDistributionEqual  ChangeDistribution = 1
	ChangeDistributionRandom ChangeDistribution = 2
)

type FeeModel interface {
	ComputeFee(tx *Transaction) (uint64, error)
}

// Fee computes the fee for the transaction.
func (tx *Transaction) Fee(f FeeModel, changeDistribution ChangeDistribution) error {
	fee, err := f.ComputeFee(tx)
	if err != nil {
		return err
	}
	satsIn := uint64(0)
	for _, i := range tx.Inputs {
		sourceSats := i.SourceTxSatoshis()
		if sourceSats == nil {
			return ErrEmptyPreviousTx
		}
		satsIn += *sourceSats
	}
	satsOut := uint64(0)
	changeOuts := uint64(0)
	for _, o := range tx.Outputs {
		if !o.Change {
			satsOut += o.Satoshis
		} else {
			changeOuts++
		}
	}
	if satsIn < satsOut+fee {
		return ErrInsufficientInputs
	}
	change := satsIn - satsOut - fee
	// There is not enough change to distribute among the change outputs.
	// We'll remove all change outputs and leave the extra for the miners.
	if changeOuts > change {
		tx.Outputs = slices.DeleteFunc(tx.Outputs, func(o *TransactionOutput) bool {
			return o.Change
		})
	} else {
		switch changeDistribution {
		case ChangeDistributionRandom:
			return errors.New("not-implemented")
		case ChangeDistributionEqual:
			changePerOutput := change / changeOuts
			for _, o := range tx.Outputs {
				if o.Change {
					o.Satoshis = changePerOutput
				}
			}
		}
	}
	return nil
}
