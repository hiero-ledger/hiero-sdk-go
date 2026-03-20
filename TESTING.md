@ -0,0 +1,94 @@
# Testing Guide

This document covers how to run tests for the Hiero Go SDK. **Note:** Running tests is only required if you are contributing to the SDK itself. If you're just using the SDK as a dependency, you don't need to run these tests.

## Unit Tests

Unit tests don't require a running Hiero network:

```bash
go test ./sdk -tags="unit" -v -timeout 9999s
```

To run a specific unit test:

```bash
go test ./sdk -tags="unit" -v -run TestYourTestName -timeout 9999s
```

## Integration (E2E) Tests

Integration tests require a running Hiero network. You can use [Solo](https://solo.hiero.org) to spin up a local development network, or connect to testnet/previewnet.

### Option 1: Using Environment Variables

```bash
export OPERATOR_ID="0.0.YOUR_ACCOUNT_ID"
export OPERATOR_KEY="YOUR_PRIVATE_KEY"
export HEDERA_NETWORK="testnet"  # or "localhost" for Solo

go test ./sdk -tags="e2e" -v -timeout 9999s
```

See [.env.sample](.env.sample) for environment variable documentation.

### Option 2: Using a Config File

```bash
export CONFIG_FILE="./path/to/config.json"

go test ./sdk -tags="e2e" -v -timeout 9999s
```

See [sdk/client-config-with-operator.json](sdk/client-config-with-operator.json) for the config file format.

### Config File Format

```json
{
  "network": {
    "0.testnet.hedera.com:50211": "0.0.3",
    "1.testnet.hedera.com:50211": "0.0.4"
  },
  "operator": {
    "accountId": "0.0.YOUR_ACCOUNT_ID",
    "privateKey": "YOUR_PRIVATE_KEY_HEX"
  }
}
```

**Note:** `HEDERA_NETWORK` takes precedence over `CONFIG_FILE` for network selection. If `HEDERA_NETWORK` is set (testnet/previewnet/localhost), `CONFIG_FILE` is ignored. When using `CONFIG_FILE`, the `OPERATOR_ID` and `OPERATOR_KEY` environment variables take precedence over the config file's operator settings.

## Local Development with Solo

[Solo](https://solo.hiero.org) is a CLI tool for running a local Hiero network. See the [Solo documentation](https://solo.hiero.org) or the [Solo repository](https://github.com/hiero-ledger/solo) for setup instructions.

Once Solo is running:

```bash
export HEDERA_NETWORK="localhost"
export OPERATOR_ID="0.0.2"
export OPERATOR_KEY="<solo-operator-key>"

go test ./sdk -tags="e2e" -v -timeout 9999s
```

## Running Examples

Examples require the same environment setup as integration tests:

```bash
# Using environment variables
export OPERATOR_ID="0.0.YOUR_ACCOUNT_ID"
export OPERATOR_KEY="YOUR_PRIVATE_KEY"

go run examples/create_account/main.go
```

Or with a config file:

```bash
export CONFIG_FILE="./path/to/config.json"

go run examples/transfer_crypto/main.go
```

