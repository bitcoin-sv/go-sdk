package transaction

// TODO: Port typescript to go
// /**
//    * Creates a new transaction, linked to its inputs and their associated merkle paths, from a BEEF (BRC-62) structure.
//    * @param beef A binary representation of a transaction in BEEF format.
//    * @returns An anchored transaction, linked to its associated inputs populated with merkle paths.
//    */
// 	 static fromBEEF(beef: number[]): Transaction {
//     const reader = new Reader(beef)
//     // Read the version
//     const version = reader.readUInt32LE()
//     if (version !== 4022206465) {
//       throw new Error(`Invalid BEEF version. Expected 4022206465, received ${version}.`)
//     }

//     // Read the BUMPs
//     const numberOfBUMPs = reader.readVarIntNum()
//     const BUMPs = []
//     for (let i = 0; i < numberOfBUMPs; i++) {
//       BUMPs.push(MerklePath.fromReader(reader))
//     }

//     // Read all transactions into an object
//     // The object has keys of TXIDs and values of objects with transactions and BUMP indexes
//     const numberOfTransactions = reader.readVarIntNum()
//     const transactions: Record<string, { pathIndex?: number, tx: Transaction }> = {}
//     let lastTXID: string
//     for (let i = 0; i < numberOfTransactions; i++) {
//       const tx = Transaction.fromReader(reader)
//       const obj: { pathIndex?: number, tx: Transaction } = { tx }
//       const txid = tx.id('hex') as string
//       if (i + 1 === numberOfTransactions) { // The last tXID is stored for later
//         lastTXID = txid
//       }
//       const hasBump = Boolean(reader.readUInt8())
//       if (hasBump) {
//         obj.pathIndex = reader.readVarIntNum()
//       }
//       transactions[txid] = obj
//     }

//     // Recursive function for adding merkle proofs or input transactions
//     const addPathOrInputs = (obj: { pathIndex?: number, tx: Transaction }): void => {
//       if (typeof obj.pathIndex === 'number') {
//         const path = BUMPs[obj.pathIndex]
//         if (typeof path !== 'object') {
//           throw new Error('Invalid merkle path index found in BEEF!')
//         }
//         obj.tx.merklePath = path
//       } else {
//         for (let i = 0; i < obj.tx.inputs.length; i++) {
//           const input = obj.tx.inputs[i]
//           const sourceObj = transactions[input.sourceTXID]
//           if (typeof sourceObj !== 'object') {
//             throw new Error(`Reference to unknown TXID in BUMP: ${input.sourceTXID}`)
//           }
//           input.sourceTransaction = sourceObj.tx
//           addPathOrInputs(sourceObj)
//         }
//       }
//     }

//	  // Read the final transaction and Add inputs and merkle proofs to the final transaction, returning it
//	  addPathOrInputs(transactions[lastTXID])
//	  return transactions[lastTXID].tx
//	}
func NewTxFromBEEF(beef []byte) *Tx {
	return nil
}

func (t *Tx) BEEF() []byte {
	return []byte{}
}
