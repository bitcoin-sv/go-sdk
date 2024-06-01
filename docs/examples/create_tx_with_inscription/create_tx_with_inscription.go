package main

import (
	"fmt"
	"log"
	"mime"
	"os"

	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template"
)

func main() {
	w, _ := wif.DecodeWIF("KznpA63DPFrmHecASyL6sFmcRgrNT9oM8Ebso8mwq1dfJF3ZgZ3V")

	tmpl := template.NewP2PKHTemplateFromPrivKey(w.PrivKey) // get public key bytes and address

	tx := transaction.NewTx()
	_ = tx.AddInputFrom(
		"39e5954ee335fdb5a1368ab9e851a954ed513f73f6e8e85eff5e31adbb5837e7",
		0,
		"76a9144bca0c466925b875875a8e1355698bdcc0b2d45d88ac",
		500,
		tmpl,
	)

	// Read the image file
	data, err := os.ReadFile("1SatLogoLight.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the content type of the image
	contentType := mime.TypeByExtension(".png")

	s, _ := tmpl.Lock()
	tx.Inscribe(&script.InscriptionArgs{
		LockingScript: s,
		Data:          data,
		ContentType:   contentType,
	})

	changeAdd, _ := script.NewAddressFromString("17ujiveRLkf2JQiGR8Sjtwb37evX7vG3WG")
	changeTmpl := template.NewP2PKHTemplateFromAddress(changeAdd)
	changeScript, _ := changeTmpl.Lock()
	tx.AddOutput(&transaction.TransactionOutput{
		LockingScript: changeScript,
		Change:        true,
	})

	err = tx.Sign()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(tx.String())
}
