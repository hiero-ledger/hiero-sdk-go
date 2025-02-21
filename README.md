![](https://img.shields.io/github/v/tag/hiero-ledger/hiero-sdk-go)
![](https://img.shields.io/github/go-mod/go-version/hiero-ledger/hiero-sdk-go)
[![](https://godoc.org/github.com/hiero-ledger/hiero-sdk-go/v2?status.svg)](http://godoc.org/github.com/hiero-project/hiero-sdk-go/v2)

# Hiero Go SDK

The Go SDK for interacting with a Hiero based network.
Hiero communicates using [gRPC](https://grpc.io);
the Protobufs definitions for the protocol are available in the [hashgraph/hedera-protobuf](https://github.com/hashgraph/hedera-protobuf) repository (the repo will be migrated to Hiero in near future).

## Usage

### Install

```sh
$ go install github.com/hiero-ledger/hiero-sdk-go/v2@latest
```

### Running Tests

# Integration
```bash
$ env CONFIG_FILE="<your_config_file>" go test ./sdk -tags="e2e" -v -timeout 9999s
```

or

```bash
$ env OPERATOR_KEY="<key>" OPERATOR_ID="<id>" go test ./sdk -tags="e2e" -timeout 9999s
```

# Unit 
```bash
$ go test ./sdk -tags="unit" -v -timeout 9999s
```

The config file _can_ contain both the network and the operator, but you can also
use environment variables `OPERATOR_KEY` and `OPERATOR_ID`. If both are provided
the network is used from the config file, but for the operator the environment variables
take precedence. If the config file is not provided then the network will default to [Hiero testnet](https://docs.hedera.com/hedera/getting-started/introduction)
and `OPERATOR_KEY` and `OPERATOR_ID` **must** be provided.

[Example Config File](./client-config-with-operator.json)

### Linting

This repository uses golangci-lint for linting. You can install a pre-commit git hook that runs golangci-lint before each commit by running the following command:

```sh
scripts/install-hooks.sh
```

## Contributing

Whether you’re fixing bugs, enhancing features, or improving documentation, your contributions are important — let’s build something great together!

Please read our [contributing guide](https://github.com/hiero-ledger/.github/blob/main/CONTRIBUTING.md) to see how you can get involved.

## Code of Conduct

Hiero uses the Linux Foundation Decentralised Trust [Code of Conduct](https://www.lfdecentralizedtrust.org/code-of-conduct).

## License

[Apache License 2.0](LICENSE)
