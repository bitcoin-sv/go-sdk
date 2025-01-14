package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/chainhash"
)

type Beef struct {
	Version      uint32
	BUMPs        []*MerklePath
	Transactions map[string]*TxOrId
}

// TODO: add methods
// makeTxidOnly
// findTxid
// findBump
// findAtomicTransaction

type DataFormat int

const (
	RawTx DataFormat = iota
	RawTxAndBumpIndex
	TxIDOnly
)

type txOrId struct {
	dataFormat  DataFormat
	KnownTxID   *chainhash.Hash
	Transaction *Transaction
}

const BEEF_V1 = uint32(4022206465) // BRC-64
const BEEF_V2 = uint32(4022206466) // BRC-96

func (t *Transaction) FromBEEF(beef []byte) error {
	tx, err := NewTransactionFromBEEF(beef)
	*t = *tx
	return err
}

func NewBEEFFromBytes(beef []byte) (*Beef, error) {
	reader := bytes.NewReader(beef)

	version, err := readVersion(reader)
	if err != nil {
		return nil, err
	}

	if version == BEEF_V1 {
		return nil, fmt.Errorf("Use NewTransactionFromBEEF to parse V1 BEEF")
	}

	BUMPs, err := readBUMPs(reader)
	if err != nil {
		return nil, err
	}

	txs, err := readTxOrIds(reader, BUMPs)
	if err != nil {
		return nil, err
	}

	return &Beef{
		Version:      version,
		BUMPs:        BUMPs,
		Transactions: txs,
	}, nil
}

func NewTransactionFromBEEF(beef []byte) (*Transaction, error) {
	reader := bytes.NewReader(beef)

	version, err := readVersion(reader)
	if err != nil {
		return nil, err
	}

	if version != BEEF_V1 {
		return nil, fmt.Errorf("Use NewBEEFFromBytes to parse anything which isn't V1 BEEF")
	}

	BUMPs, err := readBUMPs(reader)
	if err != nil {
		return nil, err
	}

	transaction, err := readTransactions(reader, BUMPs)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func readVersion(reader *bytes.Reader) (uint32, error) {
	var version uint32
	err := binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return 0, err
	}
	if version != BEEF_V1 && version != BEEF_V2 {
		return 0, fmt.Errorf("invalid BEEF version. expected %d or %d, received %d", BEEF_V1, BEEF_V2, version)
	}
	return version, nil
}

func readBUMPs(reader *bytes.Reader) ([]*MerklePath, error) {
	var numberOfBUMPs VarInt
	_, err := numberOfBUMPs.ReadFrom(reader)
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
	return BUMPs, nil
}

func readTransactions(reader *bytes.Reader, BUMPs []*MerklePath) (*Transaction, error) {
	var numberOfTransactions VarInt
	_, err := numberOfTransactions.ReadFrom(reader)
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

func readTxOrIds(reader *bytes.Reader, BUMPs []*MerklePath) (*map[string]*txOrId, error) {
	var numberOfTransactions VarInt
	_, err := numberOfTransactions.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	txs := make(map[string]*txOrId, 0)
	for i := 0; i < int(numberOfTransactions); i++ {
		formatByte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		var txOrId *txOrId
		txOrId.dataFormat = DataFormat(formatByte)

		if txOrId.dataFormat > TxIDOnly {
			return nil, fmt.Errorf("invalid data format: %d", formatByte)
		}

		if txOrId.dataFormat == TxIDOnly {
			var txid chainhash.Hash
			_, err = reader.Read(txid[:])
			txOrId.KnownTxID = &txid
			if err != nil {
				return nil, err
			}
			txs[txid.String()] = txOrId
		} else {
			bump := txOrId.dataFormat == RawTxAndBumpIndex
			// read the index of the bump
			var bumpIndex VarInt
			if bump {
				_, err := bumpIndex.ReadFrom(reader)
				if err != nil {
					return nil, err
				}
			}
			// read the transaction data
			_, err = txOrId.Transaction.ReadFrom(reader)
			if err != nil {
				return nil, err
			}
			// attach the bump
			if bump {
				txOrId.Transaction.MerklePath = BUMPs[int(bumpIndex)]
			}

			for _, input := range txOrId.Transaction.Inputs {
				sourceTxid := input.SourceTXID.String()
				if sourceObj, ok := txs[sourceTxid]; ok {
					input.SourceTransaction = sourceObj.Transaction
				} else if txOrId.Transaction.MerklePath == nil && txOrId.KnownTxID == nil {
					panic(fmt.Sprintf("Reference to unknown TXID in BUMP: %s", sourceTxid))
				}
			}

			txs[txOrId.Transaction.TxID().String()] = txOrId
		}

	}

	return &txs, nil
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
	err := binary.Write(b, binary.LittleEndian, BEEF_VERSION)
	if err != nil {
		return nil, err
	}
	bumps := []*MerklePath{}
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
		if _, ok := bumpIndex[tx.MerklePath.BlockHeight]; !ok {
			bumpIndex[tx.MerklePath.BlockHeight] = len(bumps)
			bumps = append(bumps, tx.MerklePath)
		} else {
			err := bumps[bumpIndex[tx.MerklePath.BlockHeight]].Combine(tx.MerklePath)
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
