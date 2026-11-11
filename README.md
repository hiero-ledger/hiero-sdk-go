![](https://img.shields.io/github/v/tag/hiero-ledger/hiero-sdk-go)
![](https://img.shields.io/github/go-mod/go-version/hiero-ledger/hiero-sdk-go)
[![](https://godoc.org/github.com/hiero-ledger/hiero-sdk-go/v2?status.svg)](http://godoc.org/github.com/hiero-project/hiero-sdk-go/v2)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/hiero-ledger/hiero-sdk-go/badge)](https://scorecard.dev/viewer/?uri=github.com/hiero-ledger/hiero-sdk-go)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/10697/badge)](https://bestpractices.coreinfrastructure.org/projects/10697)
[![License](https://img.shields.io/badge/license-apache2-blue.svg)](LICENSE)

# Hiero Go SDK

The Go SDK for interacting with a Hiero based network.
Hiero communicates using [gRPC](https://grpc.io);
the Protobufs definitions for the protocol are available in the [hashgraph/hedera-protobuf](https://github.com/hashgraph/hedera-protobufs) repository (the repo will be migrated to Hiero in near future).

## Prerequisites

- **Go** — A supported version of Go is required to build and run the SDK, tests, and examples. See the [go-mod/go-version](https://github.com/hiero-ledger/hiero-sdk-go/blob/main/go.mod) badge or `go.mod` for the minimum version.

## Usage

### Installation

```sh
$ go get github.com/hiero-ledger/hiero-sdk-go/v2@latest
```

### Environment Setup

To run examples or integration tests, you need Hiero network credentials. Copy [.env.sample](.env.sample) to `.env`, fill in your account details from the [Hedera Portal](https://portal.hedera.com/), then run `source .env` to export the variables.

See [TESTING.md](TESTING.md) for detailed instructions on running tests and examples.

### Review a prepared transaction body hash (candidate)

`GetTransactionBodyHash() ([]byte, error)` returns 48 raw SHA-384 bytes for a frozen transaction. For a `TransactionInterface` returned by `TransactionFromBytes`, use `TransactionGetTransactionBodyHash(transaction)`.

This is a local candidate for [SDK collaboration hub PR #288](https://github.com/hiero-ledger/sdk-collaboration-hub/pull/288). The batch framing, import ordering, rejection rules and reviewed-node execution rule are pending maintainer agreement. SDK PRs must not merge before that agreement.

Freeze with the final transaction ID, node selection and chunking before recording the digest:

```go
tx, err := hiero.NewTransferTransaction().
    SetTransactionID(hiero.TransactionIDGenerate(hiero.AccountID{Account: 1001})).
    SetNodeAccountIDs([]hiero.AccountID{{Account: 3}, {Account: 4}}).
    AddHbarTransfer(hiero.AccountID{Account: 1001}, hiero.NewHbar(-1)).
    AddHbarTransfer(hiero.AccountID{Account: 1002}, hiero.NewHbar(1)).
    Freeze()
if err != nil {
    return err
}
digest, err := tx.GetTransactionBodyHash()
if err != nil {
    return err
}
fmt.Printf("Body hash: %x\n", digest)
// Record the digest with the review, then collect transaction signatures.
```

The candidate commits to every prepared node and chunk using original body bytes. Ordinary entries use `SHA384(bodyBytes)`; the final result is `SHA384(concat(entryHashes))`, including for one entry. Batch entries use `SHA384(SHA384(outerBodyWithoutField74) || SHA384(child0.bodyBytes) || ...)`. Removing field 74 for hashing leaves every other outer-body byte intact. Outer and child envelope signatures and metadata are excluded; unknown body fields remain content.

Chunks follow transaction-ID first appearance, and every chunk follows the first chunk's node order. Node-major imports and differing later-node orders hash with that traversal without reordering stored entries. Duplicate entries, missing IDs, inconsistent node sets, malformed bodies, unsupported batch wrapper fields and nested batches return local errors. Modern `signedTransactionBytes` and complete legacy `bodyBytes` imports retain their original body bytes.

The candidate scanner validates `TransactionBody`, its transaction/node identifiers (including the transaction ID's timestamp and payer account), the atomic-batch wrapper and child `SignedTransaction` envelopes. These messages support protobuf wire types 0, 1, 2 and 5. Tags and lengths use at most five varint bytes; values use at most ten bytes without 64-bit overflow. Nonminimal encodings within those bounds remain original hash input; legacy groups are rejected. Unknown length-delimited identifier payloads remain opaque. Operation-specific submessages retain existing SDK decoding and content checks; this is not a universal protobuf validator.

Batch consistency validation uses a temporary projection of each child envelope to its original body bytes. Child signature maps, algorithm flags and unknown envelope metadata may differ between node entries; child bodies, child order and all other outer content must still match. This projection never changes stored or network envelopes.

The getter neither signs nor accesses the network. Repeated calls preserve bodies, signatures and cursors. It rejects selected IDs or nodes that duplicate or differ from the complete prepared membership before marking the body reviewed. Success locks reviewed IDs and node membership against execution regeneration or pruning while allowing later signing; failure leaves state unchanged. Reviewed imported execution resolves prepared entries by decoded transaction ID and node ID, preserving stored wire order. Unreviewed execution retains its existing cursor-driven ID regeneration and rebuilding. Unknown fields inside transaction IDs remain in the hashed and signed bytes, while identity comparisons ignore them. Existing network transaction hashes retain their signature-sensitive meaning.

In this Go binding, Freeze prepares configured common fields, including custom-fee limits and batch keys, without signing or requiring serialization before hashing. The getter rejects later pending common-field updates that differ from prepared bytes. Imported transaction IDs, batch keys and custom-fee limits retain their original nested bytes when their decoded editable values are unchanged; unknown fields remain hash input. Successful hashing marks the shared transaction body as reviewed: changing memo, fee, high-volume flag, duration, custom-fee limits or batch key is rejected using the existing frozen-mutation convention; free-function setters return that error. Indirect changes through caller-held mutable fee/key objects also cause hashing, serialization and execution to reject instead of rebuilding reviewed bytes. Re-freezing a reviewed topic transaction is rejected. Signature collection and execution options that do not change body bytes remain available. Transactions whose body hash has not been read retain their existing setter behavior.

The unit suite consumes the exact shared **35-vector** profile: **20 expected digests and 15 rejections**, including the multi-node mixed-child-signature regression. A separate five-vector identifier boundary dataset adds one digest and four rejections, for 40 shared inputs overall (21 digests and 19 rejections).

This hash identifies a prepared protobuf representation. Equivalent operations encoded independently can have different hashes; it does not promise semantic canonicalization across encoders. Re-import returned transaction bytes and recompute the digest before approval. A matching digest does not establish signature validity, authorization or network success.

### Linting

This repository uses golangci-lint for linting. You can install a pre-commit git hook that runs golangci-lint before each commit by running the following command:

```sh
scripts/install-hooks.sh
```

## Protobuf Generation Script

This script automates the process of moving and compiling `.proto` files using `protoc`. It supports multiple source directories for `proto` definitions and allows for configurable destination directory for the generation output.

### Usage

Run the script with the following command:

```sh
go run scripts/proto/generator.go -source <dir1,dir2,...> -dest <dir>
```

### Arguments

-   -source (required): A comma-separated list of directories containing .proto files.

-   -dest (required): The destination directory where all .proto files will be moved before compilation.

### Example

```sh
go run scripts/generators/proto/generator.go -source services/hapi/hedera-protobuf-java-api/src/main/proto/services/state,services/hapi/hedera-protobuf-java-api/src/main/proto/services/auxiliary,services/hapi/hedera-protobuf-java-api/src/main/proto/platform/event -dest services/hapi/hedera-protobuf-java-api/src/main/proto/services

```

### Note

-   If proto file definitions are located in the `services` submodule make sure it is initialised.
-   Keep in mind that the script does not still support protobuf import altering. If errors related to
    the `proto` import paths occur we resolve them manually.

## Request Codes and Status Generation Scripts

Scripts to automate generating the appropriate SDK counterpart implementations for statuses and codes. 

### Usage

Run the script with the following command:

```sh
go run scripts/status(or request)/generator.go -index=n
```

### Arguments

-   -index (required): The current index of the SDK statuses or request codes. Generation will be done for the statuses/codes ids `>index`.

### Example

```sh
go run scripts/status/generator.go -index=64

```

### Note

-   Keep in mind that the generated files will be located in the script directory.
-   It is advisable to always use the latest `index` as to not introduce braking changes in the generated files.
-   If further mods are required they should be done manually.

## Contributing

Whether you’re fixing bugs, enhancing features, or improving documentation, your contributions are important — let’s build something great together!

Please read our [contributing guide](https://github.com/hiero-ledger/.github/blob/main/CONTRIBUTING.md) to see how you can get involved.

## Help/Community

- Join our [community discussions](https://discord.lfdecentralizedtrust.org/) on discord.

## About Users and Maintainers

- Users and Maintainers guidelines are located in **[Hiero-Ledger's CONTRIBUTING.md file](https://github.com/hiero-ledger/.github/blob/main/CONTRIBUTING.md#about-users-and-maintainers)** under the "About-Users-and-Maintainers" section.

## Code of Conduct

Hiero uses the Linux Foundation Decentralised Trust [Code of Conduct](https://www.lfdecentralizedtrust.org/code-of-conduct).

## License

[Apache License 2.0](LICENSE)
