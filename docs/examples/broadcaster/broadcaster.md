# Example: Building a Custom Transaction Broadcast Client

This guide walks through the necessary steps for building a custom transaction broadcast client.

## Overview

A transaction broadcast client is a crucial component in any Bitcoin SV application, allowing it to communicate with the Bitcoin SV network. Implementing a transaction broadcaster can be accomplished using the clearly defined Broadcast interface.

``` go
package main

import (
        "github.com/bitcoin-sv/go-sdk/transaction"
        "github.com/bitcoin-sv/go-sdk/transaction/broadcaster"
)

func main() {

        // Create a new transaction
        hexTx := "010000000100"
        tx, _ := transaction.NewTransactionFromHex(hexTx)

        // Broadcast the transaction
        success, failure := tx.Broadcast(&broadcaster.Arc{
            ApiUrl: "https://arc.gorillapool.io",
            ApiKey: "",
        })

        // Check for errors
        if failure != nil {
            panic(failure)
        }

        // Print the success message and transaction ID
        println(success.Message, success.Txid)
}

```
