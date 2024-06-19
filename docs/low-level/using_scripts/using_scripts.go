package main

import (
	"encoding/hex"
	"fmt"
	"log"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	opcodes "github.com/bitcoin-sv/go-sdk/script"
	script "github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
)

func main() {
	// Generating and Deserializing Scripts

	// From Hex
	opTrueHex := hex.EncodeToString([]byte{opcodes.OpTRUE})
	scriptFromHex, err := script.NewFromHex(opTrueHex)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Script from Hex:", scriptFromHex)

	// From ASM
	scriptFromASM, err := script.NewFromASM("OP_DUP OP_HASH160 1451baa3aad777144a0759998a03538018dd7b4b OP_EQUALVERIFY OP_CHECKSIG")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Script from ASM:", scriptFromASM)

	// From Binary
	binaryData := []byte{opcodes.OpPUSHDATA1, 3, 1, 2, 3}
	scriptFromBinary := script.NewFromBytes(binaryData)
	fmt.Println("Script from Binary:", scriptFromBinary)

	// Advanced Example: Creating a P2PKH Locking Script
	privKey, err := ec.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	add, err := script.NewAddressFromPublicKey(privKey.PubKey(), true)
	if err != nil {
		log.Fatal(err)
	}
	// publicKeyHash := hash.Hash160(privKey.PubKey().SerialiseCompressed())
	// tmpl := template.NewP2PKHFromPubKeyEC(privKey.PubKey())
	lockingScript, err := p2pkh.Lock(add)
	if err != nil {
		log.Fatal(err)
	}
	lockingScriptASM, err := lockingScript.ToASM()
	log.Printf("Locking Script (ASM): %s\n", lockingScriptASM)
	if err != nil {
		log.Fatal(err)
	}

	// Serializing Scripts
	script, err := script.NewFromASM("OP_DUP OP_HASH160 1451baa3aad777144a0759998a03538018dd7b4b OP_EQUALVERIFY OP_CHECKSIG")
	if err != nil {
		log.Fatal(err)
	}

	// Serialize script to Hex
	scriptAsHex := script.String()
	fmt.Println("Script as Hex:", scriptAsHex)

	// Serialize script to ASM
	scriptAsASM, err := script.ToASM()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Script as ASM:", scriptAsASM)

	// Serialize script to Binary
	scriptAsBinary, err := hex.DecodeString(script.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Script as Binary:", scriptAsBinary)
}
