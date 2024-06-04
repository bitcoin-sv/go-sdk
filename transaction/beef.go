package transaction

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (t *Transaction) FromBEEF(beef []byte) error {
	t, err := NewTxFromBEEF(beef)
	return err
}

func NewTxFromBEEF(beef []byte) (*Transaction, error) {
	reader := bytes.NewReader(beef)

	var version uint32
	err := binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}
	if version != 4022206465 {
		return nil, fmt.Errorf("invalid BEEF version. expected 4022206465, received %d", version)
	}

	// Read the BUMPs
	var numberOfBUMPs VarInt
	_, err = numberOfBUMPs.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	BUMPs := make([]*MerklePath, numberOfBUMPs)
	for i := 0; i < int(numberOfBUMPs); i++ {
		BUMPs[i], err = NewMerklePathFromReader(reader)
		if err != nil {
			return nil, err
		}
	}

	// Read all transactions into an object
	var numberOfTransactions VarInt
	_, err = numberOfTransactions.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	transactions := make(map[string]*Transaction, 0)
	var tx *Transaction
	for i := 0; i < int(numberOfTransactions); i++ {
		tx = &Transaction{}
		_, err = tx.ReadFrom(reader)
		if err != nil {
			return nil, err
		}
		txid := tx.TxID()

		hasBump := make([]byte, 1)
		_, err = reader.Read(hasBump)
		if err != nil {
			return nil, err
		}
		if hasBump[0] != 0 {
			var pathIndex VarInt
			_, err = pathIndex.ReadFrom(reader)
			if err != nil {
				return nil, err
			}
			tx.MerklePath = BUMPs[int(pathIndex)]
		}
		for _, input := range tx.Inputs {
			sourceTxid := input.PreviousTxIDStr()
			if sourceObj, ok := transactions[sourceTxid]; ok {
				input.SourceTransaction = sourceObj
			} else if tx.MerklePath == nil {
				panic(fmt.Sprintf("Reference to unknown TXID in BUMP: %s", sourceTxid))
			}
		}
		transactions[txid] = tx
	}

	return tx, nil
}

func (t *Transaction) BEEF() ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, uint32(4022206465))
	bumps := map[uint32]*MerklePath{}
	bumpIndex := map[uint32]int{}
	txns := make(map[string]*Transaction, 0)
	ancestors, err := t.collectAncestors(txns)
	if err != nil {
		return nil, err
	}
	for _, txid := range ancestors {
		tx := txns[txid]
		if tx.MerklePath == nil {
			continue
		}
		if _, ok := bumps[tx.MerklePath.BlockHeight]; !ok {
			bumpIndex[tx.MerklePath.BlockHeight] = len(bumps)
			bumps[tx.MerklePath.BlockHeight] = tx.MerklePath
		} else {
			bumps[tx.MerklePath.BlockHeight].Combine(tx.MerklePath)
		}
	}

	b.Write(VarInt(len(bumps)).Bytes())
	for _, bump := range bumps {
		b.Write(bump.Bytes())
	}
	b.Write(VarInt(len(txns)).Bytes())
	for _, txid := range ancestors {
		tx := txns[txid]
		b.Write(tx.Bytes())
		if tx.MerklePath != nil {
			b.Write([]byte{1})
			b.Write(VarInt(bumpIndex[tx.MerklePath.BlockHeight]).Bytes())
		} else {
			b.Write([]byte{0})
		}
	}
	return b.Bytes(), nil
}

func (t *Transaction) collectAncestors(txns map[string]*Transaction) ([]string, error) {
	if t.MerklePath != nil {
		return []string{t.TxID()}, nil
	}
	ancestors := make([]string, 0)
	for _, input := range t.Inputs {
		if input.SourceTransaction == nil {
			return nil, fmt.Errorf("missing previous transaction for %s", t.TxID())
		}
		if _, ok := txns[input.PreviousTxIDStr()]; ok {
			continue
		}
		if grands, err := input.SourceTransaction.collectAncestors(txns); err != nil {
			return nil, err
		} else {
			ancestors = append(grands, ancestors...)
		}
	}
	ancestors = append(ancestors, t.TxID())
	return ancestors, nil
}
