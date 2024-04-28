package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"mime"
	"os"

	"github.com/bitcoin-sv/go-sdk/src/src/bscript"
	"github.com/bitcoin-sv/go-sdk/src/src/ec/wif"
	"github.com/bitcoin-sv/go-sdk/src/src/transaction"
	"github.com/bitcoin-sv/go-sdk/src/src/transaction/unlocker"
)

func main() {
	decodedWif, _ := wif.DecodeWIF("KznpA63DPFrmHecASyL6sFmcRgrNT9oM8Ebso8mwq1dfJF3ZgZ3V")

	// get public key bytes and address
	pubkey := decodedWif.SerialisePubKey()
	addr, _ := bscript.NewAddressFromPublicKeyString(hex.EncodeToString(pubkey), true)
	s, _ := bscript.NewP2PKHFromAddress(addr.AddressString)
	fmt.Println(addr.AddressString)

	tx := transaction.NewTx()

	_ = tx.From(
		"39e5954ee335fdb5a1368ab9e851a954ed513f73f6e8e85eff5e31adbb5837e7",
		0,
		"76a9144bca0c466925b875875a8e1355698bdcc0b2d45d88ac",
		500,
	)

	// Read the image file
	data, err := os.ReadFile("1SatLogoLight.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the content type of the image
	contentType := mime.TypeByExtension(".png")

	tx.Inscribe(&bscript.InscriptionArgs{
		LockingScriptPrefix: s,
		Data:                data,
		ContentType:         contentType,
	})

	err = tx.ChangeToAddress("17ujiveRLkf2JQiGR8Sjtwb37evX7vG3WG", transaction.NewFeeQuote())
	if err != nil {
		log.Fatal(err.Error())
	}

	err = tx.FillAllInputs(context.Background(), &unlocker.Getter{PrivateKey: decodedWif.PrivKey})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(tx.String())
}
