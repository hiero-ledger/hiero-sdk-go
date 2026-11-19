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

## Fuzz Tests

Fuzz tests feed generated inputs to the SDK's parsers and decoders to find inputs that make them panic. They use Go's native fuzzing and live in `sdk/*_fuzz_test.go`:

| Target | Fuzzes |
|---|---|
| `FuzzRLPItemRead` | `RLPItem.Read`, and checks that decoded items encode and decode back to the same item |
| `FuzzABIDecode` | `NewType` and `Decode` |

Each target has a seed corpus: the inputs added with `f.Add` in the test, and the files in `sdk/testdata/fuzz/<FuzzTarget>/`. The unit test command above runs every seed as a normal test, so the seed corpus is checked on every PR. New inputs are only generated with the `-fuzz` flag, which the [Fuzz Tests workflow](.github/workflows/fuzz.yml) runs every night for 5 minutes per target.

### Running a Fuzz Target

`go test -fuzz` runs one target at a time:

```bash
go test ./sdk -tags="unit" -run '^$' -fuzz '^FuzzRLPItemRead$' -fuzztime=5m
```

`-run '^$'` skips the other tests. Leave out `-fuzztime` to fuzz until you stop it with Ctrl+C. Inputs that reach new code are kept in `$(go env GOCACHE)/fuzz`, so the next run continues from them.

### Reproducing a Failure

When a target fails, Go writes the failing input to `sdk/testdata/fuzz/<FuzzTarget>/<id>` and prints how to run it again:

```bash
go test ./sdk -tags="unit" -run 'FuzzABIDecode/49098b7901fb31c9' -v
```

If the nightly workflow failed, download the `fuzz-failure-<FuzzTarget>` artifact from the failed run and extract it into `sdk/testdata/fuzz/` first.

Commit the input file together with the fix. From then on it runs as a regression test with the unit tests.

### Adding a Fuzz Target

1. Add a `Fuzz<Name>` function to a `sdk/<file>_fuzz_test.go` file with the `//go:build all || unit` build tag.
2. Seed it with `f.Add`, using the valid and invalid inputs from the matching `*_unit_test.go` file.
3. In `f.Fuzz`, call the function under test and ignore the inputs it rejects with an error. A panic fails the target on its own. Only call `t.Fatal` for a property that must hold for every accepted input, such as the round trip in `sdk/rlp_fuzz_test.go`.
4. Run the target for at least 5 minutes and fix what it finds.
5. Add the target to `matrix.target` in [.github/workflows/fuzz.yml](.github/workflows/fuzz.yml) and to the table above.

## Integration (E2E) Tests

Integration tests require a running Hiero network. You can use [Solo](https://solo.hiero.org) to spin up a local development network, or connect to testnet/previewnet.

### Option 1: Using Environment Variables

```bash
export OPERATOR_ID="0.0.YOUR_ACCOUNT_ID"
export OPERATOR_KEY="YOUR_PRIVATE_KEY"
export HEDERA_NETWORK="testnet"  # or "localhost" for Solo

go test ./sdk -tags="e2e" -v -timeout 9999s
```

See [.env.sample](.env.sample) for a template. Copy it to `.env` (gitignored) and fill in your values, then:

```bash
source .env
go test ./sdk -tags="e2e" -v -timeout 9999s
```

**Tip:** The `source .env` command only needs to be run once per terminal session. The variables remain available until you close the terminal. If you modify the `.env` file, run `source .env` again to apply the changes.

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

**Note:** `CONFIG_FILE` is only used when `OPERATOR_ID`, `OPERATOR_KEY`, and `HEDERA_NETWORK` are **all** unset. Environment variables always take precedence over config file settings.

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
export HEDERA_NETWORK="testnet"
go run examples/create_account/main.go
```

Or with a config file:

```bash
export CONFIG_FILE="./path/to/config.json"

go run examples/transfer_crypto/main.go
```
