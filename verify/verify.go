package verify

import (
	"fmt"

	"github.com/bitcoin-sv/go-sdk/script/interpreter"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/chaintracker"
)

func Verify(t *transaction.Transaction, chainTracker chaintracker.ChainTracker, feeModel transaction.FeeModel) (bool, error) {
	verifiedTxids := make(map[string]struct{})
	txQueue := []*transaction.Transaction{t}
	if chainTracker == nil {
		chainTracker = chaintracker.NewWhatsOnChain(chaintracker.MainNet, "")
	}

	for len(txQueue) > 0 {
		tx := txQueue[0]
		txQueue = txQueue[1:]
		txid := tx.TxID()
		txidStr := txid.String()

		if _, ok := verifiedTxids[txidStr]; ok {
			continue
		}

		if tx.MerklePath != nil {
			if isValid, err := tx.MerklePath.Verify(txid, chainTracker); err != nil {
				return false, err
			} else if isValid {
				verifiedTxids[txidStr] = struct{}{}
			}
		}

		if feeModel != nil {
			clone := tx.ShallowClone()
			clone.Outputs[0].Change = true
			if err := clone.Fee(feeModel, transaction.ChangeDistributionEqual); err != nil {
				return false, err
			}
			tx.TotalOutputSatoshis()
			if txFee, err := tx.GetFee(); err != nil {
				return false, err
			} else if cloneFee, err := clone.GetFee(); err != nil {
				return false, err
			} else if txFee < cloneFee {
				return false, fmt.Errorf("fee is too low")
			}
		}

		inputTotal := uint64(0)
		for vin, input := range tx.Inputs {
			if input.SourceTransaction == nil {
				return false, fmt.Errorf("input %d has no source transaction", vin)
			}
			if input.UnlockingScript == nil || len(*input.UnlockingScript) == 0 {
				return false, fmt.Errorf("input %d has no unlocking script", vin)
			}
			sourceOutput := input.SourceTransaction.Outputs[input.SourceTxOutIndex]
			inputTotal += sourceOutput.Satoshis
			sourceTxid := input.SourceTransaction.TxID().String()
			if _, ok := verifiedTxids[sourceTxid]; !ok {
				txQueue = append(txQueue, input.SourceTransaction)
			}

			otherInputs := make([]*transaction.TransactionInput, 0, len(tx.Inputs)-1)
			for i, otherInput := range tx.Inputs {
				if i != vin {
					otherInputs = append(otherInputs, otherInput)
				}
			}

			if input.SourceTXID == nil {
				input.SourceTXID = input.SourceTransaction.TxID()
			}

			if err := interpreter.NewEngine().Execute(
				interpreter.WithTx(tx, vin, sourceOutput),
				interpreter.WithForkID(),
				interpreter.WithAfterGenesis(),
			); err != nil {
				fmt.Println(err)
				return false, err
			}

		}
	}

	return true, nil
}
