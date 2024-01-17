package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/bitcoin-sv/go-sdk/crypto"
	"github.com/bitcoin-sv/go-sdk/ec"
)

// deriveMessagePoint takes a message, hashes it, and maps it to a point on the elliptic curve.
func deriveMessagePoint(message []byte) (ec.Point, string, error) {
	// Hash the message using SHA-256
	hash := crypto.Sha256(message)
	cId := hex.EncodeToString(hash[:])

	// Convert hash to a big integer
	mBn := new(big.Int).SetBytes(hash[:])

	// Get the generator point of the elliptic curve
	curve := ec.S256()
	Gx, Gy := curve.Params().Gx, curve.Params().Gy

	// Multiply the hash big integer with the generator point
	Mx, My := curve.ScalarMult(Gx, Gy, mBn.Bytes())

	return ec.Point{X: Mx, Y: My}, cId, nil
}

func main() {
	// Example usage
	message := []byte("Hello, World!")
	point, cId, err := deriveMessagePoint(message)
	if err != nil {
		panic(err)
	}

	// Output the results
	fmt.Printf("Message Point: (%s, %s)\n", point.X.String(), point.Y.String())
	fmt.Printf("cId: %s\n", cId)
}
