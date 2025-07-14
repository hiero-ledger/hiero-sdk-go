![](https://img.shields.io/github/v/tag/hiero-ledger/hiero-sdk-go)
![](https://img.shields.io/github/go-mod/go-version/hiero-ledger/hiero-sdk-go)
[![](https://godoc.org/github.com/hiero-ledger/hiero-sdk-go/v2?status.svg)](http://godoc.org/github.com/hiero-project/hiero-sdk-go/v2)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/hiero-ledger/hiero-sdk-go/badge)](https://scorecard.dev/viewer/?uri=github.com/hiero-ledger/hiero-sdk-go)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/10697/badge)](https://bestpractices.coreinfrastructure.org/projects/10697)
[![License](https://img.shields.io/badge/license-apache2-blue.svg)](LICENSE)

# Hiero Go SDK

The Go SDK for interacting with a Hiero based network.
Hiero communicates using [gRPC](https://grpc.io);
the Protobufs definitions for the protocol are available in the [hashgraph/hedera-protobuf](https://github.com/hashgraph/hedera-protobuf) repository (the repo will be migrated to Hiero in near future).

## Usage

### Installation

```sh
$ go get github.com/hiero-ledger/hiero-sdk-go/v2@latest
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

- Users and Maintainers guidelies are located in **[Hiero-Ledger's roles and groups guidelines](https://github.com/hiero-ledger/governance/blob/main/roles-and-groups.md#maintainers).**

## Code of Conduct

Hiero uses the Linux Foundation Decentralised Trust [Code of Conduct](https://www.lfdecentralizedtrust.org/code-of-conduct).

## License

[Apache License 2.0](LICENSE)
