// beef_test.go

package transaction

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	"github.com/stretchr/testify/require"
)

const BRC62Hex = "0100beef01fe636d0c0007021400fe507c0c7aa754cef1f7889d5fd395cf1f785dd7de98eed895dbedfe4e5bc70d1502ac4e164f5bc16746bb0868404292ac8318bbac3800e4aad13a014da427adce3e010b00bc4ff395efd11719b277694cface5aa50d085a0bb81f613f70313acd28cf4557010400574b2d9142b8d28b61d88e3b2c3f44d858411356b49a28a4643b6d1a6a092a5201030051a05fc84d531b5d250c23f4f886f6812f9fe3f402d61607f977b4ecd2701c19010000fd781529d58fc2523cf396a7f25440b409857e7e221766c57214b1d38c7b481f01010062f542f45ea3660f86c013ced80534cb5fd4c19d66c56e7e8c5d4bf2d40acc5e010100b121e91836fd7cd5102b654e9f72f3cf6fdbfd0b161c53a9c54b12c841126331020100000001cd4e4cac3c7b56920d1e7655e7e260d31f29d9a388d04910f1bbd72304a79029010000006b483045022100e75279a205a547c445719420aa3138bf14743e3f42618e5f86a19bde14bb95f7022064777d34776b05d816daf1699493fcdf2ef5a5ab1ad710d9c97bfb5b8f7cef3641210263e2dee22b1ddc5e11f6fab8bcd2378bdd19580d640501ea956ec0e786f93e76ffffffff013e660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac0000000001000100000001ac4e164f5bc16746bb0868404292ac8318bbac3800e4aad13a014da427adce3e000000006a47304402203a61a2e931612b4bda08d541cfb980885173b8dcf64a3471238ae7abcd368d6402204cbf24f04b9aa2256d8901f0ed97866603d2be8324c2bfb7a37bf8fc90edd5b441210263e2dee22b1ddc5e11f6fab8bcd2378bdd19580d640501ea956ec0e786f93e76ffffffff013c660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac0000000000"
const BEEF = "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"
const BEEFSet = "0200beef03fef1550d001102fd20c2009591fd79f7fb1fbd24c2fdc4911da930e1d7386f0216b6446b85eea29f978f1bfd21c202ac2a05abdae46fc2555c36a76035dedbf9fac4fc349eabffbd9d62ba440ffcb101fd116100cabeb714ea9a3f15a5e4f6138f6dd6b75bab32d8b40d178a0514e6e1e1b372f701fd8930007e04df7216a1d29bb8caabd1f78014b1b4f336eb6aee76bcf1797456ddc86b7501fd451800796afe5b113d8933f5eef2d180e72dc4b644fd76fb1243dfb791d9863702573701fd230c007a6edc003e02c429391cbf426816885731cb8054410599884eed508917a2f57c01fd100600eaa540de74506ed6abcb48e38cc544c53d373269271a7e6cf2143b7cc85d7ea401fd0903001e31aa04628b99d6cfa3e21fb4a7e773487ebc86a504e511eaff3f2176267b9401fd85010031e0d053497f85228b02879f69c4c7b43fb5abc3e0e47ea49a63853b117c9b5001c30083339d5a5b97ad77b74d3538678bb20ea7e61f8b02c24a625933eb496bebd3480160008ee445baec1613d591344a9915d77652f508e6442cd394626a3ff308bcb151f1013100f3f68f2a72e47bb41377e9e429daa496cd220bdcf702a36a209f9feba58d5552011900a01c52f4099bc7bdfea772ab03739bf009d72f24f68b5c4f8cc71a8c4da80804010d00c2ce2d5bfb9cbab9983ae1c871974f23a32c585d9b8440acc4ef5203c1d6c05401070072c7fc59a1717e90633f10d322e0f63272ae97c017d1efae04e4090abeeafac3010200a7aa5fa5576d1de6dd0e32d769592bc247be7bbd0b3e36e2d579fa1ec7d6ebce010000090cba670bea2e0d5c36e979e4cf9f79ad0874d734fb782fec2496d4c554e321010100d963646680643df73c34d7fa16f173595cf32a9ed6f64d2c8ee88a8af6b7bf52fedf590d001202fe66130200023275c6dde10d32d61af52b412b1e3956b5cd085605cd521778f11d53849fdb0cfe6713020000cd5e2298cf4d809c698c8adeeab66718e6b75b3d528bce74e6e01b984c736df901feb209010000736013454e087c89d813c99a043c9029cf2d427815c6a98ba3641c384ae52c4701fdd884007f742824bddca1582e4ded866d9609d9473397f8b86625376be74684f7fb947f01fd6d4200eb7f54ce4f920a3e4c7f96ef6b2d199c519df1b1286415581187ca608f3e47b801fd372100fa6c1c8cba3d3d5d030cd98eb91498cdffe70f0dad1000e123157d5dac22e22a01fd9a1000104c0294e478fbcac4e2325403afd86370c86043f295978b809004b2687a6c9a01fd4c08009ef5a5eaf16cab45a239c43852296ab323ca21faf256ab9768dd0a2f39970ec201fd2704006161cbd1755b66815eb69613b574920e9e836c8c3772aa2260ad3639848d520b01fd1202005e04b5afc0ea8d29dc22b611536832a2a2e7c860bbf4227ce0bdcc8a0e66284601fd0801009719f5f90e3937f3921045d202522fe315da1331acc3cce472c4b084d0debe65018500d79a1c3d45a3c41bf6526a9adbac2676159d2f3c753d7d3b6dba1dc3cbdd3c520143006b88b582d985bffc511556e471a6a20cfda2d41837245329f714214e009a3e48012000c1840dbdfc3014f1e912882b971c030fd21c0b023c01fe6fd7470d6d9bb2ab86011100f9c3de08d38588e225a5ee5334a3c03771a0b51318ca388dd1b5826951604d750109006e2b2e926c86214620d306a59522eee438a79157e9360cb76ee14a868fccc482010500d5c43ea372c432861db73ba0a6897fa29855e542a6ed910626dfb8954d94fa47010300d7863bafb5ca841ca0b13736fced1d492f0f741cb0a2beab1cafa517c878ae2c010000174ccda0879c20b85fa26d423deb0b34c5f2787127e244ccacfae39b5ba8fea7feeb590d001602fe46b3060002fa6ae8371111956f74412e3b1effcbd4fcb278124b6365b34c8cc20a5287bafffe47b306000011883eed76bdc7e7fb79efe23e3c50aa825ade46d79895de1a246e3d69a5b8cf01fea2590300009c92d7f67ac06e4bce0de4f18f438056f25138ee1a0cf61ed3a6d7f32261339b01fed0ac01000006178026214d61dc19c91cb5c08481f2f3daf03392c359de424cbd5d7135c5cf01fd69d6000174f6863438909d648fea32cdd65cbf457ab717f9be327d5d4352dbf157671e01fd356b0059536ea55010906b7071e36f78b20faaaede46a7f27ba4916dc1655836c73de701fd9b3500dee845c02c827dbcd862de359f5e6ad0ecca59213d9eb01896374d9efb7af9fd01fdcc1a00b22861b84b4537dfdaa8eb51957a51007af7836677ad14074601de6cd6c2871c01fd670d00591e76e7b07b26a6d7e940ec4f84497d9f3c7be111b15c336b24d83227db0c1001fdb20600f142d0ff9b2ddb7c21d8913f02adc7abc51fcdd5253154339450b87b59859aa601fd580300ce0307ff2027d405b8afa8a5c8834e9cc8bd073c4f463c3657562bbdb7843fe601fdad010027a3ce3a9829a3df0d9074099a6a3d76c81600a6a9c50f6cf857fb823c1a783901d700cca7689680c528f0a93fd9c980577016b37ce67ce75b1d728c4fa23008b1652b016a00b74bd3ab6c94f1216a803849afc254f37eea378c89167ff0686223db82767e3a013400434d5f48f733bb69fc5f0bd8238ffaec8d002951e6a1b52484fcc05819078372011b0053fef8153f4aed8aa8bdebeae0a6c1aa7712b84887fb565bcd9232fdd60fb0c0010c00009d9f21a9bc9e9d8c99aac9a1df47ffe02334fcb8bc8f3797d64c2564b3bf44010700838a284a4ee33c455b303e1eb23428b35d264b35c4f4b42bd6c68f1a7279f38801020042820e1ab5dbb77b0a6f266167b453f672d007d0c6eddc6229ce57c941f46c670100002c0da37e0453e7d01c810d2280a84792086b1fe1bc232e76ef6783f76c57757601010048746ad4d10a562bb53d2ed29438c9dfd0a6cacb78429277072e789d4d8dd8c101010091a52bf4a100e96dba15cbff933df60fcb26d95d6dd9b55fd5e450d5895e4526010100c202dcbdece72a45a1657ff7dbd979b031b1c8b839bc9a3b958683226644b736030100020000000140f6726035b03b90c1f770f0280444eeb041c45d026a8f4baaf00530bdc473a5020000006b483045022100ccdf467aa46d9570c4778f4e68491cc51dff4b815803d2406b6e8772d800f5ad02200ff8f11a59d207c734e9c68154dcef4023d75c37e661ab866b1d3e3ea77e6bda4121021cf99b6763736f48e6e063f99a43bfa82f15111ba0e0f9776280e6bd75d23af9ffffffff0377082800000000001976a91491b21f8856b862ff291ca0ac2ec924ba2419113788ac75330100000000001976a9144b5b285395052a61328b58c6594dd66aa6003d4988acf229f503000000001976a9148efcb6c55f5c299d48d0c74762dd811345c9093b88ac0000000001010200000001bcfe1adc5e99edb82c6a48f44cbae19bc0e5d31f9c8e4b3a92d6befb1cb2e510020000006a4730440220211655b505edd6fe9196aba77477dac5c9f638fe204243c09f1188a19164ac7f022035fb8640750515ca85df8197dec87a76db5c578f05b8ae645e30d8f70d429a324121028bf1be8161c50f98289df3ecd3185ed2273e9d448840232cf2f077f05e789c29ffffffff03d8000400000000001976a9144f427ee5f3099f0ac571f6b723a628e7b08fb64c88ac75330100000000001976a914f7cad87036406e5d3aef5d4a4d65887c76f9466788ac27db1004000000001976a9143219d1b6bd74f932dcb39a5f3b48cfde2b61cc0088ac0000000001020100000002e646efa607ff14299bc0b0cfaa65e035feb493cc440cb8abb8eb6225f8d4c1c4000000006b483045022100b410c4f82655f56fc8de4a622d3e4a8c662198de5ca8963989d70b85734986f502204fe884d99aa6ffd44bb01396b9f63bebcb7222b76e6e26c2bd60837ff555f1f8412103fda4ece7b0c9150872f8ef5241164b36a230fd9657bc43ca083d9e78bc0bcba6ffffffff3275c6dde10d32d61af52b412b1e3956b5cd085605cd521778f11d53849fdb0c000000006a473044022057f9d55ace1945866be0f83431867c58eda32d73ae3fdabed2d3424ebbe493530220553e286ae67bcaf49b0ea1d3163f41b1b3c91702a054e100c1e71ca4927f6dd8412103fda4ece7b0c9150872f8ef5241164b36a230fd9657bc43ca083d9e78bc0bcba6ffffffff04400d0300000000001976a9140e8338fa60e5391d54e99c734640e72461922d9988aca0860100000000001976a9140602787cc457f68c43581224fda6b9555aaab58e88ac10270000000000001976a91402cfbfc3931c7c1cf712574e80e75b1c2df14b2088acd5120000000000001976a914bd3dbab46060873e17ca754b0db0da4552c9a09388ac00000000"

func TestFromBEEF(t *testing.T) {
	// Decode the BEEF data from base64
	beefBytes, err := base64.StdEncoding.DecodeString(BEEF)
	require.NoError(t, err, "Failed to decode BEEF data")

	// Create a new Transaction object
	tx := &Transaction{}

	// Use the FromBEEF method to populate the transaction
	err = tx.FromBEEF(beefBytes)
	require.NoError(t, err, "FromBEEF method failed")

	expectedTxID := "ce70df889d5ba66a989b8e47294c751d19f948f004075cf265c4cbb2a7c97838"
	actualTxID := tx.TxID().String()
	require.Equal(t, expectedTxID, actualTxID, "Transaction ID does not match")
}

func TestNewBEEFFromBytes(t *testing.T) {
	// Decode the BEEF data from base64
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err, "Failed to decode BEEF data from hex string")

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err, "NewBeefFromBytes method failed")

	// Check the Beef object's properties
	require.Equal(t, uint32(4022206466), beef.Version, "Version does not match")
	require.Len(t, beef.BUMPs, 3, "BUMPs length does not match")
	require.Len(t, beef.Transactions, 3, "Transactions length does not match")
}

func TestBeefTransactionFinding(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// Test RemoveExistingTxid and findTxid
	for txid := range beef.Transactions {
		// Verify we can find it
		tx := beef.findTxid(txid)
		require.NotNil(t, tx)

		// Remove it
		beef.RemoveExistingTxid(txid)

		// Verify it's gone
		tx = beef.findTxid(txid)
		require.Nil(t, tx)
		break // just test one
	}
}

func TestBeefMakeTxidOnly(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// Get first transaction and verify it exists
	var txid string
	var originalTx *BeefTx
	for id, tx := range beef.Transactions {
		if tx.Transaction != nil {
			txid = id
			originalTx = tx
			break
		}
	}
	require.NotEmpty(t, txid)
	require.NotNil(t, originalTx)

	// Convert the hash to ensure it's valid
	hash, err := chainhash.NewHashFromHex(txid)
	require.NoError(t, err)

	// Set the KnownTxID field
	originalTx.KnownTxID = hash

	// Test MakeTxidOnly
	txidOnly := beef.MakeTxidOnly(txid)
	require.NotNil(t, txidOnly)
	require.Equal(t, TxIDOnly, txidOnly.DataFormat)
	require.NotNil(t, txidOnly.KnownTxID)
	require.Equal(t, hash.String(), txidOnly.KnownTxID.String())
}

func TestBeefSortTxs(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// First, let's check what transactions we have
	for txid, tx := range beef.Transactions {
		t.Logf("Transaction %s:", txid)
		t.Logf("  DataFormat: %v", tx.DataFormat)
		t.Logf("  Has Transaction: %v", tx.Transaction != nil)
		if tx.Transaction != nil {
			t.Logf("  Has MerklePath: %v", tx.Transaction.MerklePath != nil)
			t.Logf("  Number of Inputs: %d", len(tx.Transaction.Inputs))
		}
		t.Logf("  Has KnownTxID: %v", tx.KnownTxID != nil)
	}

	// Test SortTxs
	result := beef.SortTxs()
	require.NotNil(t, result)

	// Log the results
	t.Logf("Valid transactions: %v", result.Valid)
	t.Logf("TxIDOnly transactions: %v", result.TxidOnly)
	t.Logf("Transactions with missing inputs: %v", result.WithMissingInputs)
	t.Logf("Missing inputs: %v", result.MissingInputs)
	t.Logf("Not valid transactions: %v", result.NotValid)

	// Verify that valid transactions don't have missing inputs
	for _, txid := range result.Valid {
		require.NotContains(t, result.MissingInputs, txid, "Valid transaction should not have missing inputs")
		require.NotContains(t, result.NotValid, txid, "Valid transaction should not be in NotValid list")
		require.NotContains(t, result.WithMissingInputs, txid, "Valid transaction should not be in WithMissingInputs list")
	}

	// Verify that transactions with missing inputs are properly categorized
	for _, txid := range result.WithMissingInputs {
		require.NotContains(t, result.Valid, txid, "Transaction with missing inputs should not be in Valid list")
	}

	// Verify that invalid transactions are properly categorized
	for _, txid := range result.NotValid {
		require.NotContains(t, result.Valid, txid, "Invalid transaction should not be in Valid list")
	}
}

func TestBeefToLogString(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// Get the log string
	logStr := beef.ToLogString()

	// Verify the log string contains expected information
	require.Contains(t, logStr, "BEEF with", "Log should contain BEEF summary")
	require.Contains(t, logStr, "BUMPs", "Log should mention BUMPs")
	require.Contains(t, logStr, "Transactions", "Log should mention Transactions")
	require.Contains(t, logStr, "isValid", "Log should mention validity")

	// Verify BUMP information is logged
	require.Contains(t, logStr, "BUMP", "Log should contain BUMP details")
	require.Contains(t, logStr, "block:", "Log should contain block height")
	require.Contains(t, logStr, "txids:", "Log should contain txids")

	// Verify Transaction information is logged
	require.Contains(t, logStr, "TX", "Log should contain transaction details")
	require.Contains(t, logStr, "txid:", "Log should contain transaction IDs")

	// Verify each BUMP and transaction is mentioned
	bumpCount := beef.BUMPs
	for i := 0; i < len(bumpCount); i++ {
		require.Contains(t, logStr, fmt.Sprintf("BUMP %d", i), "Log should contain each BUMP")
	}
	for _, tx := range beef.Transactions {
		if tx.Transaction != nil {
			require.Contains(t, logStr, tx.Transaction.TxID().String(), "Log should contain each transaction ID")
		}
	}
}

func TestBeefClone(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	original, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// Clone the object
	clone := original.Clone()

	// Verify basic properties match
	require.Equal(t, original.Version, clone.Version, "Version should match")
	require.Equal(t, len(original.BUMPs), len(clone.BUMPs), "Number of BUMPs should match")
	require.Equal(t, len(original.Transactions), len(clone.Transactions), "Number of transactions should match")

	// Verify BUMPs are copied (not just referenced)
	for i, bump := range original.BUMPs {
		require.Equal(t, bump.BlockHeight, clone.BUMPs[i].BlockHeight, "BUMP BlockHeight should match")
		require.Equal(t, len(bump.Path), len(clone.BUMPs[i].Path), "BUMP Path length should match")

		// Verify each level of the path
		for j := range bump.Path {
			require.Equal(t, len(bump.Path[j]), len(clone.BUMPs[i].Path[j]), "Path level length should match")

			// Verify each PathElement
			for k := range bump.Path[j] {
				// Compare PathElement fields
				require.Equal(t, bump.Path[j][k].Offset, clone.BUMPs[i].Path[j][k].Offset, "PathElement Offset should match")
				if bump.Path[j][k].Hash != nil {
					require.Equal(t, bump.Path[j][k].Hash.String(), clone.BUMPs[i].Path[j][k].Hash.String(), "PathElement Hash should match")
				}
				if bump.Path[j][k].Txid != nil {
					require.Equal(t, *bump.Path[j][k].Txid, *clone.BUMPs[i].Path[j][k].Txid, "PathElement Txid should match")
				}
				if bump.Path[j][k].Duplicate != nil {
					require.Equal(t, *bump.Path[j][k].Duplicate, *clone.BUMPs[i].Path[j][k].Duplicate, "PathElement Duplicate should match")
				}
			}
		}
	}

	// Verify transactions are copied (not just referenced)
	for txid, tx := range original.Transactions {
		clonedTx, exists := clone.Transactions[txid]
		require.True(t, exists, "Transaction should exist in clone")
		require.Equal(t, tx.DataFormat, clonedTx.DataFormat, "Transaction DataFormat should match")
		if tx.Transaction != nil {
			require.Equal(t, tx.Transaction.TxID().String(), clonedTx.Transaction.TxID().String(), "Transaction ID should match")
		}
		if tx.KnownTxID != nil {
			require.Equal(t, tx.KnownTxID.String(), clonedTx.KnownTxID.String(), "KnownTxID should match")
		}
	}

	// Modify clone and verify original is unchanged
	clone.Version = 999
	require.NotEqual(t, original.Version, clone.Version, "Modifying clone should not affect original")

	// Remove a transaction from clone and verify original is unchanged
	for txid := range clone.Transactions {
		delete(clone.Transactions, txid)
		_, exists := original.Transactions[txid]
		require.True(t, exists, "Removing transaction from clone should not affect original")
		break // just test one
	}
}

func TestBeefTrimknownTxIDs(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// Convert some transactions to TxIDOnly format
	var txidsToTrim []string
	for txid, tx := range beef.Transactions {
		if tx.Transaction != nil {
			// Convert to TxIDOnly and add to our list to trim
			beef.MakeTxidOnly(txid)
			txidsToTrim = append(txidsToTrim, txid)
			if len(txidsToTrim) >= 2 { // Convert 2 transactions to test with
				break
			}
		}
	}
	require.GreaterOrEqual(t, len(txidsToTrim), 1, "Should have at least one transaction to trim")

	// Verify the transactions are now in TxIDOnly format
	for _, txid := range txidsToTrim {
		tx := beef.findTxid(txid)
		require.NotNil(t, tx)
		require.Equal(t, TxIDOnly, tx.DataFormat)
	}

	// Trim the known TxIDs
	beef.TrimknownTxIDs(txidsToTrim)

	// Verify the transactions were removed
	for _, txid := range txidsToTrim {
		tx := beef.findTxid(txid)
		require.Nil(t, tx, "Transaction should have been removed")
	}

	// Verify other transactions still exist
	for txid, tx := range beef.Transactions {
		require.NotContains(t, txidsToTrim, txid, "Remaining transaction should not have been in trim list")
		if tx.DataFormat == TxIDOnly {
			require.NotContains(t, txidsToTrim, txid, "TxIDOnly transaction that wasn't in trim list should still exist")
		}
	}
}

func TestBeefGetValidTxids(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// First, let's check what transactions we have
	t.Log("Checking transactions in BEEF:")
	for txid, tx := range beef.Transactions {
		t.Logf("Transaction %s:", txid)
		t.Logf("  DataFormat: %v", tx.DataFormat)
		t.Logf("  Has Transaction: %v", tx.Transaction != nil)
		if tx.Transaction != nil {
			t.Logf("  Has MerklePath: %v", tx.Transaction.MerklePath != nil)
			t.Logf("  Number of Inputs: %d", len(tx.Transaction.Inputs))
			for i, input := range tx.Transaction.Inputs {
				t.Logf("    Input %d SourceTXID: %s", i, input.SourceTXID.String())
			}
		}
		t.Logf("  Has KnownTxID: %v", tx.KnownTxID != nil)
	}

	// Get sorted transactions to see what's valid
	sorted := beef.SortTxs()
	t.Log("\nSorted transaction results:")
	t.Logf("  Valid: %v", sorted.Valid)
	t.Logf("  TxidOnly: %v", sorted.TxidOnly)
	t.Logf("  WithMissingInputs: %v", sorted.WithMissingInputs)
	t.Logf("  MissingInputs: %v", sorted.MissingInputs)
	t.Logf("  NotValid: %v", sorted.NotValid)

	// Get valid txids
	validTxids := beef.GetValidTxids()
	t.Logf("\nGetValidTxids result: %v", validTxids)

	// Verify results match
	require.Equal(t, sorted.Valid, validTxids, "GetValidTxids should return same txids as SortTxs.Valid")

	// If we have any valid transactions, verify they exist and have valid inputs
	if len(validTxids) > 0 {
		for _, txid := range validTxids {
			tx := beef.findTxid(txid)
			require.NotNil(t, tx, "Valid txid should exist in transactions map")

			// If it has a transaction, verify it has no missing inputs
			if tx.Transaction != nil {
				for _, input := range tx.Transaction.Inputs {
					sourceTx := beef.findTxid(input.SourceTXID.String())
					require.NotNil(t, sourceTx, "Input transaction should exist for valid transaction")
				}
			}
		}
	} else {
		t.Log("No valid transactions found - this is expected if all transactions have missing inputs or are not valid")
	}
}

func TestBeefFindTransactionForSigning(t *testing.T) {
	// Decode the BEEF data from hex string
	beefBytes, err := hex.DecodeString(BEEFSet)
	require.NoError(t, err)

	// Create a new Beef object
	beef, err := NewBeefFromBytes(beefBytes)
	require.NoError(t, err)

	// First, let's check what transactions we have
	t.Log("Checking transactions in BEEF:")
	for txid, tx := range beef.Transactions {
		t.Logf("Transaction %s:", txid)
		t.Logf("  DataFormat: %v", tx.DataFormat)
		t.Logf("  Has Transaction: %v", tx.Transaction != nil)
		if tx.Transaction != nil {
			t.Logf("  Has MerklePath: %v", tx.Transaction.MerklePath != nil)
			t.Logf("  Number of Inputs: %d", len(tx.Transaction.Inputs))
			for i, input := range tx.Transaction.Inputs {
				t.Logf("    Input %d SourceTXID: %s", i, input.SourceTXID.String())
			}
		}
		t.Logf("  Has KnownTxID: %v", tx.KnownTxID != nil)
	}

	// Get sorted transactions to see what's valid
	sorted := beef.SortTxs()
	t.Log("\nSorted transaction results:")
	t.Logf("  Valid: %v", sorted.Valid)
	t.Logf("  TxidOnly: %v", sorted.TxidOnly)
	t.Logf("  WithMissingInputs: %v", sorted.WithMissingInputs)
	t.Logf("  MissingInputs: %v", sorted.MissingInputs)
	t.Logf("  NotValid: %v", sorted.NotValid)

	// Get valid txids
	validTxids := beef.GetValidTxids()
	t.Logf("\nGetValidTxids result: %v", validTxids)

	// For this test, we'll use any transaction that has full data
	var testTxid string
	for txid, tx := range beef.Transactions {
		if tx.Transaction != nil {
			testTxid = txid
			break
		}
	}
	require.NotEmpty(t, testTxid, "Should have at least one transaction with full data")

	// Test FindTransactionForSigning
	tx := beef.FindTransactionForSigning(testTxid)
	require.NotNil(t, tx, "Should find a transaction for signing")
	require.Equal(t, testTxid, tx.TxID().String(), "Transaction ID should match")
}
