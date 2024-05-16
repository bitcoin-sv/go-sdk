# Example: Creating a Simple Transaction
This guide walks you through the steps of creating a simple Bitcoin transaction. To get started, let's explain some basic concepts around Bitcoin transactions.

## Understanding and Creating Transactions
Transactions in Bitcoin are mechanisms for transferring value and invoking smart contract logic. The Transaction class in the BSV SDK encapsulates the creation, signing, and broadcasting of transactions, also enabling the use of Bitcoin's scripting language for locking and unlocking coins.

## Creating and Signing a Transaction
Consider the scenario where you need to create a transaction. The process involves specifying inputs (where the bitcoins are coming from) and outputs (where they're going). Here's a simplified example:

```
package main

import (
	"context"
	"log"

	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/unlocker"
)

// https://goplay.tools/snippet/bnsS-pA56ob
func main() {
	// Create a new transaction
	tx := transaction.NewTx()

	// Add the inputs
	_ = tx.From(
		// Previous transaction ID (hex)
		"11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d",
		// Previous transaction output index
		0,
		// Previous transaction script (hex)
		"76a9144bca0c466925b875875a8e1355698bdcc0b2d45d88ac",
		// Previous transaction output value in satoshis
		1500,
	)

	// Add the outputs
	_ = tx.PayToAddress(
		// Destination address
		"1NRoySJ9Lvby6DuE2UQYnyT67AASwNZxGb",
		// Value in satoshis
		1000,
	)

	// Fill all inputs with the given private key
	decodedWif, _ := wif.DecodeWIF("KznvCNc6Yf4iztSThoMH6oHWzH9EgjfodKxmeuUGPq5DEX5maspS")

	// Sign the transaction
	if err := tx.FillAllInputs(
		// The default context which is empty but can be 
		// used to control cancellations and timeouts.
		context.Background(), 
		&unlocker.Getter{PrivateKey: decodedWif.PrivKey}); err != nil {
			log.Fatal(err.Error())
	}
	log.Printf("tx: %s\n", tx)
}

```

### Package and Imports:

The program is defined in the main package.
It imports necessary packages for handling the transaction, including context management, logging, WIF (Wallet Import Format) decoding, and transaction-related functionalities from the go-sdk.
Main Function:

`transaction.NewTx()` creates a new, empty transaction.

`tx.From(...)` adds an input to the transaction. This input includes:

- The transaction ID (txid) from which the input is being spent.
- The output index (vout) in the referenced transaction.
- A script (locking script) representing the locking conditions under which the output can be spent.
- The number of satoshis from the output that is being spent.

`tx.PayToAddress(...)` adds an output to the transaction, specifying:

- The recipient's Bitcoin SV address.
- The amount of satoshis to send to that address.

`wif.DecodeWIF(...)` decodes a WIF-encoded private key. This is necessary for signing the transaction inputs.

`tx.FillAllInputs(...)` attempts to sign all inputs of the transaction using the private key provided. This is essential for validating the transaction on the network.

The context passed `(context.Background())` is the default context which is empty but can be used to control cancellations and timeouts.
Error Handling and Logging:

If `tx.FillAllInputs(...)` returns an error, the program logs the error and terminates using `log.Fatal(...)`.

If no error occurs, it logs the serialized form of the transaction, which is ready to be broadcasted to the network.