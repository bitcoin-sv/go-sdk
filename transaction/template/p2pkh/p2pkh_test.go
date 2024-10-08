package p2pkh_test

import (
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/require"
)

func TestLocalUnlocker_UnlockAllInputs(t *testing.T) {
	t.Parallel()

	incompleteTx := "010000000193a35408b6068499e0d5abd799d3e827d9bfe70c9b75ebe209c91d25072326510000000000ffffffff02404b4c00000000001976a91404ff367be719efa79d76e4416ffb072cd53b208888acde94a905000000001976a91404d03f746652cfcb6cb55119ab473a045137d26588ac00000000"
	tx, err := transaction.NewTransactionFromHex(incompleteTx)
	require.NoError(t, err)
	require.NotNil(t, tx)

	prevTx := transaction.NewTransaction()
	prevTx.Outputs = make([]*transaction.TransactionOutput, tx.InputIdx(0).SourceTxOutIndex+1)
	prevTx.Outputs[tx.InputIdx(0).SourceTxOutIndex] = &transaction.TransactionOutput{Satoshis: 100000000}
	prevTx.Outputs[tx.InputIdx(0).SourceTxOutIndex].LockingScript, err = script.NewFromHex("76a914c0a3c167a28cabb9fbb495affa0761e6e74ac60d88ac")
	require.NoError(t, err)
	tx.Inputs[0].SourceTransaction = prevTx

	// Our private key
	priv, err := ec.PrivateKeyFromWif("cNGwGSc7KRrTmdLUZ54fiSXWbhLNDc2Eg5zNucgQxyQCzuQ5YRDq")
	require.NoError(t, err)

	unlocker, err := p2pkh.Unlock(priv, nil)
	require.NoError(t, err)

	s, err := unlocker.Sign(tx, 0)
	require.NoError(t, err)
	tx.Inputs[0].UnlockingScript = s

	expectedSignedTx := "010000000193a35408b6068499e0d5abd799d3e827d9bfe70c9b75ebe209c91d2507232651000000006b483045022100c1d77036dc6cd1f3fa1214b0688391ab7f7a16cd31ea4e5a1f7a415ef167df820220751aced6d24649fa235132f1e6969e163b9400f80043a72879237dab4a1190ad412103b8b40a84123121d260f5c109bc5a46ec819c2e4002e5ba08638783bfb4e01435ffffffff02404b4c00000000001976a91404ff367be719efa79d76e4416ffb072cd53b208888acde94a905000000001976a91404d03f746652cfcb6cb55119ab473a045137d26588ac00000000"
	require.Equal(t, expectedSignedTx, tx.String())
	require.NotEqual(t, incompleteTx, tx.String())
}

func TestLocalUnlocker_ValidSignature(t *testing.T) {
	tests := map[string]struct {
		tx *transaction.Transaction
	}{
		"valid signature 1": {
			tx: func() *transaction.Transaction {
				tx := transaction.NewTransaction()
				require.NoError(t, tx.AddInputFrom("45be95d2f2c64e99518ffbbce03fb15a7758f20ee5eecf0df07938d977add71d", 0, "76a914c7c6987b6e2345a6b138e3384141520a0fbc18c588ac", 15564838601, nil))

				script1, err := script.NewFromHex("76a91442f9682260509ac80722b1963aec8a896593d16688ac")
				require.NoError(t, err)

				tx.AddOutput(&transaction.TransactionOutput{
					Satoshis:      375041432,
					LockingScript: script1,
				})

				script2, _ := script.NewFromHex("76a914c36538e91213a8100dcb2aed456ade363de8483f88ac")
				tx.AddOutput(&transaction.TransactionOutput{
					Satoshis:      15189796941,
					LockingScript: script2,
				})

				return tx
			}(),
		},
		"valid signature 2": {
			tx: func() *transaction.Transaction {
				tx := transaction.NewTransaction()

				require.NoError(
					t,
					tx.AddInputFrom("64faeaa2e3cbadaf82d8fa8c7ded508cb043c5d101671f43c084be2ac6163148", 1, "76a914343cadc47d08a14ef773d70b3b2a90870b67b3ad88ac", 5000000000, nil),
				)
				tx.Inputs[0].SequenceNumber = 0xfffffffe

				script1, err := script.NewFromHex("76a9140108b364bbbddb222e2d0fac1ad4f6f86b10317688ac")
				require.NoError(t, err)
				tx.AddOutput(&transaction.TransactionOutput{
					Satoshis:      2200000000,
					LockingScript: script1,
				})

				script2, err := script.NewFromHex("76a9143ac52294c730e7a4e9671abe3e7093d8834126ed88ac")
				require.NoError(t, err)
				tx.AddOutput(&transaction.TransactionOutput{
					Satoshis:      2799998870,
					LockingScript: script2,
				})
				return tx
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tx := test.tx

			priv, err := ec.PrivateKeyFromWif("cNGwGSc7KRrTmdLUZ54fiSXWbhLNDc2Eg5zNucgQxyQCzuQ5YRDq")
			require.NoError(t, err)

			unlocker, err := p2pkh.Unlock(priv, nil)
			require.NoError(t, err)
			uscript, err := unlocker.Sign(tx, 0)
			require.NoError(t, err)

			tx.Inputs[0].UnlockingScript = uscript
			parts, err := script.DecodeScript(*tx.Inputs[0].UnlockingScript)
			require.NoError(t, err)

			sigBytes := parts[0].Data
			publicKeyBytes := parts[1].Data

			publicKey, err := ec.ParsePubKey(publicKeyBytes)
			require.NoError(t, err)

			sig, err := ec.ParseDERSignature(sigBytes)
			require.NoError(t, err)

			sh, err := tx.CalcInputSignatureHash(0, sighash.AllForkID)
			require.NoError(t, err)

			require.True(t, sig.Verify(sh, publicKey))
		})
	}
}

// type mockUnlockerGetter struct {
// 	t            *testing.T
// 	unlockerFunc func(ctx context.Context, lockingScript *script.Script) (transaction.ScriptTemplate, error)
// }

// func (m *mockUnlockerGetter) Unlocker(ctx context.Context, lockingScript *script.Script) (transaction.ScriptTemplate, error) {
// 	require.NotNil(m.t, m.unlockerFunc, "unlockerFunc not set in this test")
// 	return m.unlockerFunc(ctx, lockingScript)
// }

// type mockUnlocker struct {
// 	t      *testing.T
// 	script string
// }

// func (m *mockUnlocker) UnlockingScript(ctx context.Context, tx *transaction.Transaction, params transaction.UnlockParams) (*script.Script, error) {
// 	uscript, err := script.NewFromASM(m.script)
// 	require.NoError(m.t, err)

// 	return uscript, nil
// }

// func TestLocalUnlocker_NonSignature(t *testing.T) {
// 	t.Parallel()
// 	tests := map[string]struct {
// 		tx                  *transaction.Transaction
// 		unlockerFunc        func(ctx context.Context, lockingScript *script.Script) (transaction.Unlocker, error)
// 		expUnlockingScripts []string
// 	}{
// 		"simple script": {
// 			tx: func() *transaction.Transaction {
// 				tx := transaction.NewTransaction()
// 				require.NoError(t, tx.AddInputFrom("45be95d2f2c64e99518ffbbce03fb15a7758f20ee5eecf0df07938d977add71d", 0, "52529387", 15564838601))
// 				return tx
// 			}(),
// 			unlockerFunc: func(ctx context.Context, lockingScript *script.Script) (transaction.Unlocker, error) {
// 				asm, err := lockingScript.ToASM()
// 				require.NoError(t, err)

// 				unlocker, ok := map[string]*mockUnlocker{
// 					"OP_2 OP_2 OP_ADD OP_EQUAL": {t: t, script: "OP_4"},
// 				}[asm]

// 				require.True(t, ok)
// 				require.NotNil(t, unlocker)

// 				return unlocker, nil
// 			},
// 			expUnlockingScripts: []string{"OP_4"},
// 		},
// 		"multiple inputs unlocked": {
// 			tx: func() *transaction.Transaction {
// 				tx := transaction.NewTransaction()
// 				require.NoError(t, tx.AddInputFrom("45be95d2f2c64e99518ffbbce03fb15a7758f20ee5eecf0df07938d977add71d", 0, "52529487", 15564838601))
// 				require.NoError(t, tx.AddInputFrom("45be95d2f2c64e99518ffbbce03fb15a7758f20ee5eecf0df07938d977add71d", 0, "52589587", 15564838601))
// 				require.NoError(t, tx.AddInputFrom("45be95d2f2c64e99518ffbbce03fb15a7758f20ee5eecf0df07938d977add71d", 0, "5a559687", 15564838601))
// 				return tx
// 			}(),
// 			unlockerFunc: func(ctx context.Context, lockingScript *script.Script) (transaction.ScriptTemplate, error) {
// 				asm, err := lockingScript.ToASM()
// 				require.NoError(t, err)

// 				unlocker, ok := map[string]*mockUnlocker{
// 					"OP_2 OP_2 OP_SUB OP_EQUAL":  {t: t, script: "OP_FALSE"},
// 					"OP_2 OP_8 OP_MUL OP_EQUAL":  {t: t, script: "OP_16"},
// 					"OP_10 OP_5 OP_DIV OP_EQUAL": {t: t, script: "OP_2"},
// 				}[asm]

// 				require.True(t, ok)
// 				require.NotNil(t, unlocker)

// 				return unlocker, nil
// 			},
// 			expUnlockingScripts: []string{"OP_FALSE", "OP_16", "OP_2"},
// 		},
// 	}

// 	for name, test := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			tx := test.tx
// 			require.Equal(t, len(tx.Inputs), len(test.expUnlockingScripts))

// 			ug := &mockUnlockerGetter{
// 				t:            t,
// 				unlockerFunc: test.unlockerFunc,
// 			}
// 			require.NoError(t, tx.FillAllInputs(context.Background(), ug))
// 			for i, script := range test.expUnlockingScripts {
// 				asm, err := tx.Inputs[i].UnlockingScript.ToASM()
// 				require.NoError(t, err)

// 				require.Equal(t, script, asm)
// 			}
// 		})
// 	}
// }

//
// func TestBareMultiSigValidation(t *testing.T) {
// 	txHex := "0100000001cfb38c76cadeb5b96c3863d9e298fe96e24e594b75f69c37aa709f45b76d1b25000000009200483045022100d83dc84d3ea3fb36b006f6887e1e16811c59fe9a9b79b84142874a90d5b834160220052967be98c26270de0082b0fecab5a40d5bc48d5034b6cdfc2b8e47210e1469414730440220099ffa89363f9a05f23a4fa318ddbefeeeec4b41f6abde7083a3be6696ed904902201722110a488df3780a260ba09b7de6363bfce7f6beec9819e9b9f47f6e978d8141ffffffff01a8840100000000001976a91432b996f742e774b0241be9007f831558ba06d20b88ac00000000"
// 	tx, err := transaction.NewTransactionFromHex(txHex)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	// txid := tx.GetTxID()
// 	// fmt.Println(txid)
//
// 	var sigs = make([]*ec.Signature, 2)
// 	var sigHashTypes = make([]uint32, 2)
// 	var publicKeys = make([]*ec.PublicKey, 3)
//
// 	sigScript := tx.Inputs[0].UnlockingScript
//
// 	sig0Bytes := []byte(*sigScript)[2:73]
// 	sig0HashType, _ := binary.Uvarint([]byte(*sigScript)[73:74])
// 	sig1Bytes := []byte(*sigScript)[75:145]
// 	sig1HashType, _ := binary.Uvarint([]byte(*sigScript)[145:146])
//
// 	pk0, _ := hex.DecodeString("023ff15e2676e03b2c0af30fc17b7fb354bbfa9f549812da945194d3407dc0969b")
// 	pk1, _ := hex.DecodeString("039281958c651c013f5b3b007c78be231eeb37f130b925ceff63dc3ac8886f22a3")
// 	pk2, _ := hex.DecodeString("03ac76121ffc9db556b0ce1da978021bd6cb4a5f9553c14f785e15f0e202139e3e")
//
// 	publicKeys[0], err = ec.ParsePubKey(pk0, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	publicKeys[1], err = ec.ParsePubKey(pk1, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	publicKeys[2], err = ec.ParsePubKey(pk2, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	sigs[0], err = ec.ParseDERSignature(sig0Bytes, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	sigs[1], err = ec.ParseDERSignature(sig1Bytes, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	sigHashTypes[0] = uint32(sig0HashType)
// 	sigHashTypes[1] = uint32(sig1HashType)
//
// 	var previousTxSatoshis uint64 = 99728
// 	var SourceTxScript, _ = script.NewFromHex("5221023ff15e2676e03b2c0af30fc17b7fb354bbfa9f549812da945194d3407dc0969b21039281958c651c013f5b3b007c78be231eeb37f130b925ceff63dc3ac8886f22a32103ac76121ffc9db556b0ce1da978021bd6cb4a5f9553c14f785e15f0e202139e3e53ae")
// 	var prevIndex uint32
// 	var outIndex uint32
//
// 	for i, sig := range sigs {
// 		sighash := signature.GetSighashForInputValidation(tx, sigHashTypes[i], outIndex, prevIndex, previousTxSatoshis, SourceTxScript)
// 		h, err := hex.DecodeString(sighash)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		for j, pk := range publicKeys {
// 			valid := sig.Verify(util.ReverseBytes(h), pk)
// 			t.Logf("signature %d against pulbic key %d => %v\n", i, j, valid)
// 		}
//
// 	}
//
// }
//
// func TestP2SHMultiSigValidation(t *testing.T) { // NOT working properly!
// 	txHex := "0100000001d0219010e1f74ec8dd264a63ef01b5c72aab49a74c9bffd464c7f7f2b193b34700000000fdfd0000483045022100c2ffae14c7cfae5c1b45776f4b2d497b0d10a9e3be55b1386c555f90acd022af022025d5d1d33429fabd60c41763f9cda5c4b64adbddbd90023febc005be431b97b641473044022013f65e41abd6be856e7c7dd7527edc65231e027c42e8db7358759fc9ccd77b7d02206e024137ee54d2fac9f1dce858a85cb03fb7ba93b8e015d82e8a959b631f91ac414c695221021db57ae3de17143cb6c314fb206b56956e8ed45e2f1cbad3947411228b8d17f1210308b00cf7dfbb64604475e8b18e8450ac6ec04655cfa5c6d4d8a0f3f141ee419421030c7f9342ff6583599db8ee8b52383cadb4cf6fee3650c1ad8f66158a4ff0ebd953aefeffffff01b70f0000000000001976a91415067448220971206e6b4d90733d70fe9610631688ac56750900"
// 	tx, err := transaction.NewTransactionFromHex(txHex)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	// txid := tx.GetTxID()
// 	// fmt.Println(txid)
//
// 	var sigs = make([]*ec.Signature, 2)
// 	var sigHashTypes = make([]uint32, 2)
// 	var publicKeys = make([]*ec.PublicKey, 3)
//
// 	sigScript := tx.Inputs[0].UnlockingScript
//
// 	sig0Bytes := []byte(*sigScript)[2:73]
// 	sig0HashType, _ := binary.Uvarint([]byte(*sigScript)[73:74])
// 	sig1Bytes := []byte(*sigScript)[75:145]
// 	sig1HashType, _ := binary.Uvarint([]byte(*sigScript)[145:146])
//
// 	pk0, _ := hex.DecodeString("021db57ae3de17143cb6c314fb206b56956e8ed45e2f1cbad3947411228b8d17f1")
// 	pk1, _ := hex.DecodeString("0308b00cf7dfbb64604475e8b18e8450ac6ec04655cfa5c6d4d8a0f3f141ee4194")
// 	pk2, _ := hex.DecodeString("030c7f9342ff6583599db8ee8b52383cadb4cf6fee3650c1ad8f66158a4ff0ebd9")
//
// 	publicKeys[0], err = ec.ParsePubKey(pk0, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	publicKeys[1], err = ec.ParsePubKey(pk1, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	publicKeys[2], err = ec.ParsePubKey(pk2, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	sigs[0], err = ec.ParseDERSignature(sig0Bytes, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	sigs[1], err = ec.ParseDERSignature(sig1Bytes, ec.S256())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	sigHashTypes[0] = uint32(sig0HashType)
// 	sigHashTypes[1] = uint32(sig1HashType)
//
// 	var previousTxSatoshis uint64 = 8785040
// 	var SourceTxScript, _ = script.NewFromHex("5221021db57ae3de17143cb6c314fb206b56956e8ed45e2f1cbad3947411228b8d17f1210308b00cf7dfbb64604475e8b18e8450ac6ec04655cfa5c6d4d8a0f3f141ee419421030c7f9342ff6583599db8ee8b52383cadb4cf6fee3650c1ad8f66158a4ff0ebd953ae")
// 	var prevIndex uint32 = 1
// 	var outIndex uint32 = 0
//
// 	for i, sig := range sigs {
// 		sighash := signature.GetSighashForInputValidation(tx, sigHashTypes[i], outIndex, prevIndex, previousTxSatoshis, SourceTxScript)
// 		h, err := hex.DecodeString(sighash)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		for j, pk := range publicKeys {
// 			valid := sig.Verify(util.ReverseBytes(h), pk)
// 			t.Logf("signature %d against pulbic key %d => %v\n", i, j, valid)
// 		}
//
// 	}
//
// }
