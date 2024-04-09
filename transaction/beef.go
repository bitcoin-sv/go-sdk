package transaction

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BeefTx struct {
	pathIndex *uint64
	tx        *Tx
}

func NewTxFromBEEF(beef []byte) (*Tx, error) {
	reader := bytes.NewReader(beef)

	var version uint32
	err := binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}
	if version != 4022206465 {
		return nil, fmt.Errorf("Invalid BEEF version. Expected 4022206465, received %d", version)
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

	transactions := make(map[string]*BeefTx, 0)
	lastTxid := ""
	for i := 0; i < int(numberOfTransactions); i++ {
		tx := &Tx{}
		_, err = tx.ReadFrom(reader)
		if err != nil {
			return nil, err
		}
		beefTx := &BeefTx{tx: tx}
		txid := tx.TxID()
		if i+1 == int(numberOfTransactions) {
			lastTxid = txid
		}
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
			val := uint64(pathIndex)
			beefTx.pathIndex = &val
		}
		transactions[txid] = beefTx
	}

	populateInputsFromBeef(transactions[lastTxid], BUMPs, transactions)
	return transactions[lastTxid].tx, nil
}

func populateInputsFromBeef(beefTx *BeefTx, bumps []*MerklePath, transactions map[string]*BeefTx) {
	if beefTx.pathIndex != nil {
		path := bumps[*beefTx.pathIndex]
		if path == nil {
			panic("Invalid merkle path index found in BEEF!")
		}
		beefTx.tx.MerklePath = path
	} else {
		for _, input := range beefTx.tx.Inputs {
			sourceTxid := input.PreviousTxIDStr()
			if sourceObj, ok := transactions[sourceTxid]; !ok {
				panic(fmt.Sprintf("Reference to unknown TXID in BUMP: %s", sourceTxid))
			} else {
				input.PreviousTxScript = sourceObj.tx.Outputs[input.PreviousTxOutIndex].LockingScript
				input.PreviousTxSatoshis = sourceObj.tx.Outputs[input.PreviousTxOutIndex].Satoshis
				populateInputsFromBeef(sourceObj, bumps, transactions)
			}
		}
	}
}

func (t *Tx) BEEF() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, uint32(4022206465))
	bumps := make([]*MerklePath, 0)
	txs := make(map[string]*BeefTx, 0)

	addPathsAndInputs(t, &bumps, txs)
	b.Write(VarInt(len(bumps)).Bytes())
	for _, bump := range bumps {
		b.Write(bump.Bytes())
	}
	b.Write(VarInt(len(txs)).Bytes())
	for _, beefTx := range txs {
		b.Write(beefTx.tx.Bytes())
		if beefTx.pathIndex != nil {
			b.Write([]byte{1})
			b.Write(VarInt(*beefTx.pathIndex).Bytes())
		} else {
			b.Write([]byte{0})
		}
	}
	return b.Bytes()
}

func addPathsAndInputs(tx *Tx, bumps *[]*MerklePath, txs map[string]*BeefTx) {
	beefTx := &BeefTx{tx: tx}
	hasProof := tx.MerklePath != nil

	if hasProof {
		added := false
		for i, bump := range *bumps {
			pathIndex := uint64(i)
			if bytes.Equal(bump.Bytes(), tx.MerklePath.Bytes()) {
				beefTx.pathIndex = &pathIndex
				added = true
				break
			}
			if bump.BlockHeight == tx.MerklePath.BlockHeight {
				rootA, _ := bump.ComputeRoot(nil)
				rootB, _ := tx.MerklePath.ComputeRoot(nil)
				if rootA == rootB {
					bump.Combine(tx.MerklePath)
					beefTx.pathIndex = &pathIndex
					added = true
					break
				}
			}
		}
		if !added {
			pathIndex := uint64(len(*bumps))
			beefTx.pathIndex = &pathIndex
			*bumps = append(*bumps, tx.MerklePath)
		}
	}
}
