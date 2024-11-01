# BSV BLOCKCHAIN | Software Development Kit for Go

Welcome to the BSV Blockchain Libraries Project, the comprehensive Go SDK designed to provide an updated and unified layer for developing scalable applications on the BSV Blockchain. This SDK addresses the limitations of previous tools by offering a fresh, peer-to-peer approach, adhering to SPV, and ensuring privacy and scalability.

# Status
![golangci-lint](https://github.com/bitcoin-sv/go-sdk/workflows/golangci-lint/badge.svg)

## Table of Contents

- [BSV BLOCKCHAIN | Software Development Kit for Go](#bsv-blockchain--software-development-kit-for-go)
  - [Table of Contents](#table-of-contents)
  - [Objective](#objective)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Basic Usage](#basic-usage)
    - [Examples](#examples)
  - [Features \& Deliverables](#features--deliverables)
  - [Documentation](#documentation)
  - [Contribution Guidelines](#contribution-guidelines)
  - [Support \& Contacts](#support--contacts)
  - [License](#license)

## Objective

The BSV Blockchain Libraries Project aims to structure and maintain a middleware layer of the BSV Blockchain technology stack. By facilitating the development and maintenance of core libraries, it serves as an essential toolkit for developers looking to build on the BSV Blockchain.

## Getting Started

### Installation

To install the SDK, run:

```bash
go install github.com/bitcoin-sv/go-sdk
```

### Basic Usage

Here's a [simple example](https://goplay.tools/snippet/bnsS-pA56ob) of using the SDK to create and sign a transaction:

```go
package main

import (
	"context"
	"log"

	wif "github.com/bitcoin-sv/go-sdk/compat/wif"
	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/unlocker"
)

func main() {
	tx := transaction.NewTransaction()

	unlockingScriptTemplate, _ := p2pkh.Unlock(priv, nil)
   if err := tx.AddInputFrom(
		"11b476ad8e0a48fcd40807a111a050af51114877e09283bfa7f3505081a1819d",
		0,
		"76a9144bca0c466925b875875a8e1355698bdcc0b2d45d88ac",
		1500,
		unlockingScriptTemplate,
	); err!= nil {
      log.Fatal(err.Error())
   }

	_ = tx.PayToAddress("1AdZmoAQUw4XCsCihukoHMvNWXcsd8jDN6", 1000)

	priv, _ := ec.PrivateKeyFromWif("KznvCNc6Yf4iztSThoMH6oHWzH9EgjfodKxmeuUGPq5DEX5maspS")

	if err := tx.Sign(); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("tx: %s\n", tx)
}

```

See the [Go Doc](https://pkg.go.dev/github.com/bitcoin-sv/go-sdk) for a complete list of available modules and functions.

### Examples

Check out the [examples folder](https://github.com/bitcoin-sv/go-sdk/tree/master/docs/examples) for more advanced examples.

## Features

- **Performance Oriented**: Designed to deliver performant functionality for large scale / high demand systems.
- **Cryptographic Primitives**: Secure key management, signature computations, and encryption protocols.
- **Script Level Constructs**: Network-compliant script interpreter with support for custom scripts and serialization formats.
- **Transaction Construction and Signing**: Comprehensive transaction builder API, ensuring versatile and secure transaction creation.
- **Transaction Broadcast Management**: Mechanisms to send transactions to both miners and overlays, ensuring extensibility and future-proofing.
- **Merkle Proof Verification**: Tools for representing and verifying merkle proofs, adhering to various serialization standards.
- **Serializable SPV Structures**: Structures and interfaces for full SPV verification.
- **Secure Encryption and Signed Messages**: Enhanced mechanisms for encryption and digital signatures, replacing outdated methods.
- **Shamir Key Splitting & Recombining**: Allows private keys to be split into N shares, and recombined by providing M of N shares.
- **Compatability Packages**: Supports additional / deprecated features like ECIES, Bitcoin Signed Message, and BIP32 style key derivation.

## Documentation

The SDK is richly documented with code-level annotations. This should show up well within editors like VSCode. For complete API docs, check out the [Go Doc](https://pkg.go.dev/github.com/bitcoin-sv/go-sdk). Please refer to the [Libraries Wiki](#) (link to be provided) for a deep dive into each feature, tutorials, and usage examples.

## Contribution Guidelines

We're always looking for contributors to help us improve the SDK. Whether it's bug reports, feature requests, or pull requests - all contributions are welcome.

1. **Fork & Clone**: Fork this repository and clone it to your local machine.
2. **Set Up**: Run `go get github.com/bitcoin-sv/go-sdk` to get all the modules.
3. **Make Changes**: Create a new branch and make your changes.
4. **Test**: Ensure all tests pass by running `go test ./...`.
5. **Commit**: Commit your changes and push to your fork.
6. **Pull Request**: Open a pull request from your fork to this repository.

For more details, check the [contribution guidelines](./CONTRIBUTING.md).

For information on past releases, check out the [changelog](./CHANGELOG.md). For future plans, check the [roadmap](./ROADMAP.md)!

## Support & Contacts

Project Owners: Thomas Giacomo and Darren Kellenschwiler

Development Team Lead: Luke Rohenaz

For questions, bug reports, or feature requests, please open an issue on GitHub or contact us directly.

## License

The license for the code in this repository is the Open BSV License. Refer to [LICENSE.txt](./LICENSE.txt) for the license text.

Thank you for being a part of the BSV Blockchain Libraries Project. Let's build the future of BSV Blockchain together!
