package transaction_test

import (
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/assert"
)

const BRC62Hex = "0100beef01fe636d0c0007021400fe507c0c7aa754cef1f7889d5fd395cf1f785dd7de98eed895dbedfe4e5bc70d1502ac4e164f5bc16746bb0868404292ac8318bbac3800e4aad13a014da427adce3e010b00bc4ff395efd11719b277694cface5aa50d085a0bb81f613f70313acd28cf4557010400574b2d9142b8d28b61d88e3b2c3f44d858411356b49a28a4643b6d1a6a092a5201030051a05fc84d531b5d250c23f4f886f6812f9fe3f402d61607f977b4ecd2701c19010000fd781529d58fc2523cf396a7f25440b409857e7e221766c57214b1d38c7b481f01010062f542f45ea3660f86c013ced80534cb5fd4c19d66c56e7e8c5d4bf2d40acc5e010100b121e91836fd7cd5102b654e9f72f3cf6fdbfd0b161c53a9c54b12c841126331020100000001cd4e4cac3c7b56920d1e7655e7e260d31f29d9a388d04910f1bbd72304a79029010000006b483045022100e75279a205a547c445719420aa3138bf14743e3f42618e5f86a19bde14bb95f7022064777d34776b05d816daf1699493fcdf2ef5a5ab1ad710d9c97bfb5b8f7cef3641210263e2dee22b1ddc5e11f6fab8bcd2378bdd19580d640501ea956ec0e786f93e76ffffffff013e660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac0000000001000100000001ac4e164f5bc16746bb0868404292ac8318bbac3800e4aad13a014da427adce3e000000006a47304402203a61a2e931612b4bda08d541cfb980885173b8dcf64a3471238ae7abcd368d6402204cbf24f04b9aa2256d8901f0ed97866603d2be8324c2bfb7a37bf8fc90edd5b441210263e2dee22b1ddc5e11f6fab8bcd2378bdd19580d640501ea956ec0e786f93e76ffffffff013c660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac0000000000"

func TestNewTransaction(t *testing.T) {
	t.Parallel()
	t.Run("create tx", func(t *testing.T) {
		// examlpe wif and associated address
		priv, _ := ec.PrivateKeyFromWif("KznvCNc6Yf4iztSThoMH6oHWzH9EgjfodKxmeuUGPq5DEX5maspS")
		address, _ := script.NewAddressFromPublicKey(priv.PubKey(), true)

		// Source transaction data
		sourceRawtx := "010000000138c7c61c14ffb063c3bb2664041a3e29ea6ea0412a0c18ff725ba4e9e12afae2030000006a47304402203e9ab8e4c14addf3b4741540b556cfb0e0efb67dc1a7b5ce84c3ac56b3fd447802203c9f49f7bd893ebd7060176dfc36bcaff9d2c443d9a0dd6cd2d59b372c024d20412102798913bc057b344de675dac34faafe3dc2f312c758cd9068209f810877306d66ffffffff02dc050000000000002076a914eb0bd5edba389198e73f8efabddfc61666969ff788ac6a0568656c6c6faa0d0000000000001976a914eb0bd5edba389198e73f8efabddfc61666969ff788ac00000000"
		sourceMerklePathHex := "fed7c509000a02fddd010069464172a5d0cd3d641516166091ab84d230e8848ac9ccdc93f185d7b1b07902fddc01029d81a1815050f3a7bf8392e077481151af50a011a10708d4fc480a8ead76b41101ef00a93658c713530e49e2d6cad2529ecf06eb20620b9e1d3bdf75dbef8f509a5cc101760040808d97bfcb804293013e2108c4df25996ea9ba517671ff721c7be73dbfc3c5013a000435fef874132a7ebda11760ad63eccf37ba82f41793d6453f744b0873829c77011c000a0d32242d744e2007e8c3ccbfd761380d7c4340a90d8255cd608ad307752cd8010f007718f8c034a5ff0adf9c3c337660c4592bd6a6ff10de2d8f01afbb8c65f9143e0106006214d394450c84eabdcf04e7ecc6b893e1649ecc48bb3a6f38d48afcb0f2bc6a0102006a5fe10c65d3ce6950b4cbbd2bd584bcec0263c5178c3226bde14d7e307d4557010000a14df02e34b74d15dbcd0c7896b3dfb8ffb136cc3ba61ec118b37ddc70974cd5010100a5f147afb93db1ffe573b69b7c84abc3582c6cd7f3eaf82b4142d7557c28f0ae"

		sourceTransaction, _ := transaction.NewTransactionFromHex(sourceRawtx)
		merklePath, _ := transaction.NewMerklePathFromHex(sourceMerklePathHex)
		sourceTransaction.MerklePath = merklePath

		// Create a new transaction
		tx := transaction.NewTransaction()

		// Create a new P2PKH unlocker from the private key
		unlocker, err := p2pkh.Unlock(priv, nil)
		assert.NoError(t, err)
		// Add an input
		tx.AddInputFromTx(sourceTransaction, 0, unlocker)

		// Add the outputs
		tx.AddOutput(&transaction.TransactionOutput{
			LockingScript: p2pkh.Lock(address),
			Satoshis:      1,
		})
		assert.NoError(t, err)

		// Sign the transaction
		if err := tx.Sign(); err != nil {
			assert.NoError(t, err)
		}

		_, err = tx.BEEF()
		assert.NoError(t, err)
	})
}

func TestBEEF(t *testing.T) {
	t.Parallel()
	t.Run("deserialize and serialize", func(t *testing.T) {
		tx, err := transaction.NewTransactionFromBEEFHex(BRC62Hex)
		assert.NoError(t, err)
		assert.Equal(t, tx.Inputs[0].SourceTransaction.MerklePath.BlockHeight, uint32(814435))
		beef, err := tx.BEEFHex()
		assert.NoError(t, err)
		assert.Equal(t, beef, BRC62Hex)
	})
}

func TestEF(t *testing.T) {
	t.Run("Serialization and deserialization", func(t *testing.T) {
		tx, err := transaction.NewTransactionFromBEEFHex(BRC62Hex)
		assert.NoError(t, err)
		ef, err := tx.EFHex()
		assert.NoError(t, err)
		assert.Equal(t, ef, "010000000000000000ef01ac4e164f5bc16746bb0868404292ac8318bbac3800e4aad13a014da427adce3e000000006a47304402203a61a2e931612b4bda08d541cfb980885173b8dcf64a3471238ae7abcd368d6402204cbf24f04b9aa2256d8901f0ed97866603d2be8324c2bfb7a37bf8fc90edd5b441210263e2dee22b1ddc5e11f6fab8bcd2378bdd19580d640501ea956ec0e786f93e76ffffffff3e660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac013c660000000000001976a9146bfd5c7fbe21529d45803dbcf0c87dd3c71efbc288ac00000000")
	})
}
