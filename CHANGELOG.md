# CHANGELOG

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Table of Contents

- [Unreleased](#unreleased)
- [1.0.0 - YYYY-MM-DD](#100---yyyy-mm-dd)

## [Unreleased]
- porting in all optimizations by Teranode team to their active go-bt fork
- introducing chainhash to remove type coercion on tx hashes through the project
- remove ByteStringLE (replaced by chainhash)
- update opRshift and opLshift modeled after C code in node software and tested against failing vectors
- add tests and vectors for txs using opRshift that were previously failing to verify
- update examples
- lint - change international spellings to match codebase standards, use require instead of assert, etc
- add additional test vectors from known failing transactions

### Added
- `MerkePath.ComputeRootHex`
- `MerklePath.VerifyHex`

### Changed
- `SourceTXID` on `TransactionInput` is now a `ChainHash` instead of `[]byte`
- `IsValidRootForHeight` on `ChainTracker` now takes a `ChainHash` instead of `[]byte`
- `MerklePath.ComputeRoot` now takes a `ChainHash` instead of a hex `string`.
- `MerklePath.Verify` now takes a `ChainHash` instead of hex `string`.
- `Transaction.TxID` now returns a `ChainHash` instead of a hex `string`
- `Transaction.PreviousOutHash` was renamed to `SourceOutHash`, and returns a `ChainHash` instead of `[]byte`
- The `TxID` field of the `UTXO` struct in the `transaction` package is now a `ChainHash` instead of `[]byte`
- Renamed `TransactionInput.SetPrevTxFromOutput` to `SetSourceTxFromOutput`

### Removed
- `TransactionInput.PreviousTxIDStr`
- `Transaction.TxIDBytes`
- `UTXO.TxIDStr` in favor of `UTXO.TxID.String()`

### Fixed
- `opcodeRShift` and `opcodeLShift` was fixed to match node logic and properly execute scripts using `OP_RSHIFT` and `OP_LSHIFT`.

---

## [1.0.0] - YYYY-MM-DD

### Added
- Initial release

---

### Template for New Releases:

Replace `X.X.X` with the new version number and `YYYY-MM-DD` with the release date:

```
## [X.X.X] - YYYY-MM-DD

### Added
- 

### Changed
- 

### Deprecated
- 

### Removed
- 

### Fixed
- 

### Security
- 
```

Use this template as the starting point for each new version. Always update the "Unreleased" section with changes as they're implemented, and then move them under the new version header when that version is released.