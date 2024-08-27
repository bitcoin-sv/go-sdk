package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func (t *Transaction) FromBEEF(beef []byte) error {
	tx, err := NewTransactionFromBEEF(beef)
	*t = *tx
	return err
}

func NewTransactionFromBEEF(beef []byte) (*Transaction, error) {
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
			sourceTxid := input.SourceTXID.String()
			if sourceObj, ok := transactions[sourceTxid]; ok {
				input.SourceTransaction = sourceObj
			} else if tx.MerklePath == nil {
				panic(fmt.Sprintf("Reference to unknown TXID in BUMP: %s", sourceTxid))
			}
		}
		transactions[txid.String()] = tx
	}

	return tx, nil
}

func NewTransactionFromBEEFHex(beefHex string) (*Transaction, error) {
	if beef, err := hex.DecodeString(beefHex); err != nil {
		return nil, err
	} else {
		return NewTransactionFromBEEF(beef)
	}
}

func (t *Transaction) BEEF() ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, uint32(4022206465))
	if err != nil {
		return nil, err
	}
	bumps := map[uint32]*MerklePath{}
	bumpIndex := map[uint32]int{}
	txns := map[string]*Transaction{t.TxID().String(): t}
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
			err := bumps[tx.MerklePath.BlockHeight].Combine(tx.MerklePath)
			if err != nil {
				return nil, err
			}
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

func (t *Transaction) BEEFHex() (string, error) {
	if beef, err := t.BEEF(); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(beef), nil
	}
}

func (t *Transaction) collectAncestors(txns map[string]*Transaction) ([]string, error) {
	if t.MerklePath != nil {
		return []string{t.TxID().String()}, nil
	}
	ancestors := make([]string, 0)
	for _, input := range t.Inputs {
		if input.SourceTransaction == nil {
			return nil, fmt.Errorf("missing previous transaction for %s", t.TxID())
		}
		if _, ok := txns[input.SourceTXID.String()]; ok {
			continue
		}
		txns[input.SourceTXID.String()] = input.SourceTransaction
		if grands, err := input.SourceTransaction.collectAncestors(txns); err != nil {
			return nil, err
		} else {
			ancestors = append(grands, ancestors...)
		}
	}
	ancestors = append(ancestors, t.TxID().String())
	return ancestors, nil
}
