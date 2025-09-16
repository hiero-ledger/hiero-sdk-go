## v2.70.0

### Added
- Validation to `NodeUpdateTransaction` setters to align with JS SDK [#1490](https://github.com/hiero-ledger/hiero-sdk-go/pull/1490)
- `slices.Equal` to simplify code [#1476](https://github.com/hiero-ledger/hiero-sdk-go/pull/1476) 

### Fixed
- Status codes `StatusServiceEndpointsExceededLimit`, `StatusInvalidGossipCaCertificate`, `StatusInvalidIPV4Address` has swapped places [#1486](https://github.com/hiero-ledger/hiero-sdk-go/pull/1486)
- Freeze error not being checked in some flows [#1486](https://github.com/hiero-ledger/hiero-sdk-go/pull/1486)

## v2.69.0

### Added
- The ability to set AutoRenewPeriod in seconds
- Makes SetStakedAccountID and SetStakedNodeID explicitly mutually exclusive
- Support for Private and Public keys in SetAdminKey for contract update
- Adds the ability to remove autorenew account via setting it to `0.0.0`
- `customFeeLimits` for Scheduled Transactions (currently only for `TopicSubmitTransaction`)

## v2.68.0

### Added
- Introduced new method `DeleteGrpcWebProxyEndpoint` for `NodeUpdateTransaction` that effectively sets the proxy to `null` in the mirror node api and removes it from the address book state. [#1444](https://github.com/hiero-ledger/hiero-sdk-go/pull/1444)
- Chunk size option for `TopicMessageSubmitTransaction`, added `SetChunkSize` and `GetChunkSize` [#1448](https://github.com/hiero-ledger/hiero-sdk-go/pull/1448)

### Changed
- Now it's possible to update the topic memo to empty string using the setter, as previously the value would be ignored [#1448](https://github.com/hiero-ledger/hiero-sdk-go/pull/1448)
- Migrated to Go version 1.24 [#1452](https://github.com/hiero-ledger/hiero-sdk-go/pull/1452)

## v2.67.0

### Added
- Checks for `NodeUpdate` and `NodeDelete` to Return error when `nodeID` is not explicitly set [#1426](https://github.com/hiero-ledger/hiero-sdk-go/pull/1426)
    - This prevents users from accidentally updating or deleting node with id 0

### Changed
- Refactored usage of for loops with `slices.Contains`, which can make the code more concise and easy to read [#1442](https://github.com/hiero-ledger/hiero-sdk-go/pull/1442)

## v2.66.0

### Added
- `ToEvmAddress` method in AccountId, ContractId, DelegateContractId, FileId, TokenId, and TopicId [#1391](https://github.com/hiero-ledger/hiero-sdk-go/pull/1391)
These methods remove encoding shards and realm from evm address methods according to this proposal https://github.com/hiero-ledger/sdk-collaboration-hub/blob/main/proposals/remove-shard-and-realm-encoding-from-evm-address-generation.md

### Deprecated
- `FromSolidityAddress` in AccountId, ContractId, DelegateContractId, FileId, TokenId, and TopicId [#1391](https://github.com/hiero-ledger/hiero-sdk-go/pull/1391)
- `ToSolidityAddress` in AccountId, ContractId, DelegateContractId, FileId, TokenId, and TopicId [#1391](https://github.com/hiero-ledger/hiero-sdk-go/pull/1391)
- `EthereumFlow`, with the introduction of jumbo transactions, it should always be less cost and more efficient to use `EthereumTransaction` instead [#1428](https://github.com/hiero-ledger/hiero-sdk-go/pull/1428)

## v2.65.0

### Added
- Add persistent shard and realm support to Client [#1395](https://github.com/hiero-ledger/hiero-sdk-go/pull/1395)
  - `ClientForNetworkV2` : extracts the shard and realm from the network map and persists it
  - `ClientFromConfig` : can now accept shard and realm into the json schema
  - `GetShard`
  - `GetRealm`

### Deprecated 
- `ClientForMirrorNetworkWithRealmAndShard` : Use `ClientForMirrorNetworkWithShardAndRealm` instead.
- `ClientForNetwork` : Use `ClientForNetworkV2` instead.

### Changed 
- Replaced bip39 library (https://github.com/tyler-smith/go-bip39), since it is no longer available [#1418](https://github.com/hiero-ledger/hiero-sdk-go/pull/1418)

## v2.64.0

### Added
- Transaction size APIs [#1392](https://github.com/hiero-ledger/hiero-sdk-go/pull/1392)
  - `GetTransactionSize`
  - `GetTransactionBodySize`
  - `GetTransactionBodySizeAllChunks`

### Fixed
- Pause/UnpauseTransaction's protobuf methods [#1393](https://github.com/hiero-ledger/hiero-sdk-go/pull/1393)

## v2.63.0

### Changed
- Auto-setting `autoRenewAccount` for `TokenCreateTransaction` only if autorenew period is specified. [#1386](https://github.com/hiero-ledger/hiero-sdk-go/pull/1384)

### Fixed
- Now we validate `FileAppendTransaction`'s receipts for all chunks. This improves UX and removes false-positives. [#1379](https://github.com/hiero-ledger/hiero-sdk-go/pull/1379)

### Added
- HIP-1064: Daily Rewards For Active Nodes https://hips.hedera.com/hip/hip-1064 [#1383](https://github.com/hiero-ledger/hiero-sdk-go/pull/1383)
- HIP-1046: gRPC-Web proxy endpoints to the Address Book https://hips.hedera.com/hip/hip-1046 [#1383](https://github.com/hiero-ledger/hiero-sdk-go/pull/1383) 
New APIs for `NodeCreate` and `NodeUpdate`:
    - Endpoint GetGrpcWebProxyEndpoint()
    - SetGrpcWebProxyEndpoint(Endpoint)
    - bool GetDeclineReward()
    - SetDeclineReward(bool)

- Offline multi-node signing support [#1378](https://github.com/hiero-ledger/hiero-sdk-go/pull/1378)    
New APIs for `Transaction.go`
  - `GetSignableBodyBytes` : returns a list of SignableBody objects for each signed transaction in the transaction list.
  - `AddSignatureV2` : adds signature for multi-node and multi-chunk transactions.

## v2.62.0

### Added
- New APIs `GetData` and `SetData` in `EthereumTransactionData`. Used to modify the call data in ethereum [#1336](https://github.com/hiero-ledger/hiero-sdk-go/pull/1336)
- Response struct, containing the `NodeID` is returned by `EthereumFlow.Execute` even if the transaction fails [#1336](https://github.com/hiero-ledger/hiero-sdk-go/pull/1336)

### Fixed
- Fixed an issue where the receipt children were not included when the transaction throttled, when `SetIncludeChildren` is set to true [#1368](https://github.com/hiero-ledger/hiero-sdk-go/pull/1368)

## v2.61.0

### Added
- New APIs for handling of non-zero shard and realms for static files [#1363](https://github.com/hiero-ledger/hiero-sdk-go/pull/1363)
  - FileId.getAddressBookFileIdFor(uint64 realm, uint64 shard)
  - FileId.getFeeScheduleFileIdFor(uint64 realm, uint64 shard)
  - FileId.getExchangeRatesFileIdFor(uint64 realm, uint64 shard)
  - Client.forMirrorNetwork(List<string>, uint64 realm, uint64 shard)
- Support for HIP-551 Batch Transaction https://hips.hedera.com/hip/hip-551
  It defines a mechanism to execute batch transactions such that a series of transactions (HAPI calls) depending on each other can be rolled into one transaction that passes the ACID test (atomicity, consistency, isolation, and durability). [#1347](https://github.com/hiero-ledger/hiero-sdk-go/pull/1347)
    - New BatchTransaction struct that has a list of innerTransactions and innerTransactionIds.
    - New `batchKey` field in Transaction class that must sign the BatchTransaction
    - New `batchify` method that sets the batch key and marks a transaction as part of a batch transaction (inner transaction). The transaction is signed by the client of the operator and frozen.
- Extend `SetKeyWithAlias` funcs to support `PublicKey` [#1348](https://github.com/hiero-ledger/hiero-sdk-go/pull/1348)
- Support for deserializing transaction bytes, representing single transaction proto body. [#1347](https://github.com/hiero-ledger/hiero-sdk-go/pull/1347)

## v2.60.0

### Added
- Support for HIP-1021: Improve Assignment of Auto-Renew Account ID for Topics (https://hips.hedera.com/hip/hip-1021). The autoRenewAccountId will automatically be set to the payer_account_id of the transaction
  if an Admin Key is not provided during topic creation [#1355](https://github.com/hiero-ledger/hiero-sdk-go/pull/1355)
- Added a User-Agent header to outgoing gRPC requests via a unary interceptor. The header value includes the SDK identifier (hiero-sdk-go) and the version obtained from build information (defaulting to DEV if unavailable). This aids in tracking SDK usage metrics [#1315](https://github.com/hiero-ledger/hiero-sdk-go/pull/1315)

### Fixed
-  Fixed `INVALID_NODE_ACCOUNT` error when setting nodes for paid queries.
The issue was that we generated payment transactions for all the nodes in the q.nodeAccountIds list and not the current node we are pointing to. This caused problems because q.nodeAccountIDs._Advance() is called when we get the cost for the query and this moves the pointer to the next node in the list.

## v2.59.0

### Added
- `EIP-2930` transaction type compatibility. A new struct `EthereumEIP2930Transaction` is added for RLP encoding `EIP-2930` transactions. [#1325](https://github.com/hiero-ledger/hiero-sdk-go/pull/1325)
- Specifying min TLS version for gRPC communication. [#1308](https://github.com/hiero-ledger/hiero-sdk-go/pull/1308)
- PublicKey `VerifySignedMessage` method in place of `Verify`. [#1314](https://github.com/hiero-ledger/hiero-sdk-go/pull/1314)
- `PrivateKey.GetRecoveryId` method. This method retrieves the recovery ID (also known as the 'v' value) associated with ECDSA signatures,
facilitating signature verification processes. [#1324](https://github.com/hiero-ledger/hiero-sdk-go/pull/1324)

### Changed
- Modification of the `PrivateKey.Sign` method output. The `Sign` method for ECDSA private keys has been updated to return only the r and s components of the signature,
reducing the output from 65 bytes to 64 bytes.This change aligns the SDK's behavior with standard ECDSA signature formats, which typically include only the r and s values. [#1324](https://github.com/hiero-ledger/hiero-sdk-go/pull/1324)

### Deprecated
- PublicKey `Verify` since it's not keytype agnostic and has different behavior for ed25519 and ecdsa keys. #1314
- `VerifySignature` is no longer maintained, since it requires the full 65 byte signature and pre-hashing using Keccak256Hash. `PublicKey.VerifySignedMessage` is preferred. [#1324](https://github.com/hiero-ledger/hiero-sdk-go/pull/1324)

### Fixed
- The PublicKey `VerifyTransaction` method was building the proto transaction body, which overrides the signatures and causes `INVALID_SIGNATURE` error. 
The build logic is now removed and a new check if the pubkey is in the transaction was added. [#1314](https://github.com/hiero-ledger/hiero-sdk-go/pull/1314)

## v2.58.0

### Fixed
- `TokenPauseTransaction` and `TokenUnpauseTransaction` fromBytes.

### Removed
- Automatic setting of autorenew account for topic create.

### Added
- `ScheduledNetworkUpdate` example.

## v2.57.0

### Added
- Docs for min/max values of some parameters.
- Support for HIP-1021: improve assignment of auto renew account id for topics

### Fixed
- Errors are returning as `nil` in 2 functions for `ContractFunctionParameters`.
- `EthereumFlow` for creating large contracts.

## v2.56.0

### Removed
- `AccountStakersQuery`, since it was not supported in consensus node for a long time and now it's removed permanently.

### Deprecated
- `Livehash` transactions and queries `SystemDeleteTransaction` and `SystemUndeleteTransaction`.

### Changed
- `NftId` string format from `serial@tokenid` to `tokenid@serial`.

### Added
- Support for HIP-991: revenue generating topics.

### Fixed
- Keeping the `transactionFee` for a transaction while serializing/deserializing when the transaction was not frozen.

## v2.55.0

### Added

- New APIs in `AccountCreateTransaction` : `SetECDSAKeyWithAlias(ECDSAKey)`, `setKeyWithAlias(Key, ECDSAKey)` and `setKeyWithoutAlias(Key)`.

### Changed

- Deprecated `setKey` in `AccountCreateTransaction`.
- Set default max fee for transaction to 2 HBars.

## v2.54.0

### Added 
- `DeleteTokenNftAllowanceAllSerials()` method in `AccountAllowanceApproveTransaction`

### Fixed
- Overriding the default values of properties in some transactions when doing `transaction.toBytes()` and then `TransactionFromBytes(transaction)`.

### Changed
- Moved all source files to `/sdk` directory. The new way of importing the SDK is `import hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"`

## v2.53.0

### Added
- `NextExchangeRate` in the `TransactionReceipt`.
- `ScheduleRef` in the `TransactionRecord`.
- 2 new queries `MirrorNodeContractCallQuery` and `MirrorNodeContractEstimateGasQuery` for estimation/simulation of contract operations.
- Missing `RequestType`s in `request_type.go`.
- Validation for creating ECDSA Public keys from bytes.

### Changed
- Removed `ExpectedDecimals` from the `TransactionRecord`, as 'dead' property.

## v2.52.0

### Added
- Support for Long Term Scheduled Transactions (HIP-423).

### Changed
- Renamed packages to hiero

## v2.51.0

### Fixed
- Max backoff for grpc requests
- Resubmit transaction in case of throttle status at record

### Added
- Getters for internal keys and threshold for `KeyList`
- New api for creating client without schedule network update - `ClientFromConfigWithoutScheduleNetworkUpdate`
- New api for creating client with mirror network - `ClientForMirrorNetwork`

## v2.50.0

### Changed

-   Implemented generics in `Transaction.go`
-   Replaced go-ethereum library

## v2.49.0

### Changed

-   Replace `/common/math` package from go-ethereum
-   Update protobufs from `hedera-services`

## v2.48.0

### Fixed

-   Reset `stakedAccountID` when setting `stakedNodeID` and vice versa
-   Fix `FEE_SCHEDULE_FILE_PART_UPLOADED` marked as error

## v2.47.0

### Added

-   Functionality to pass a string in `SetMessage` function for `TopicMessageSubmitTransaction`

### Fixed

-   Resubmit transaction in case of throttle status at receipt

## v2.46.0

### Added

-   `TokenClaimAirdropTransaction` and `TokenCancelAirdropTransaction` (part of HIP-904)

### Fixed

-   Handling of `0x` prefix when constructing ECDSA keys

## v2.45.0

### Added

-   `TokenAirdropTransaction` (part of HIP-904)

### Fixed

-   Handling of `THROTTLED_AT_CONSENSUS` status code

## v2.44.0

### Added

-   `NodeCreateTransaction`,`NodeUpdateTransaction`,`NodeDeleteTransaction` (part of HIP-869)

## v2.43.0

### Added

-   `Key` functions such as `KeyFromBytes` `KeyToBytes`
-   `KeyList` functions such as `SetThreshold`

## v2.42.0

### Added

-   `TokenReject` functionality (part of HIP-904)

### Fixed

-   `TransactionReceiptQuery` and `AccountBalanceQuery` execution flows

## v2.41.0

### Added

-   Modified `AccountUpdateTransaction` to allow `maxAutomaticTokenAssociations` to support `-1` as a valid value

## v2.40.0

### Added

-   Implemented custom derivation paths in Menmonic ECDSA private key derivation

### Fixed

-   Revisited and fix failing examples
-   Gracefully handle `PlatformNotActive` status code

## v2.39.0

### Added

-   Implemented HIP-540: Change or remove existing keys from a token

## v2.38.1

### Changed

-   `AccountBalanceQuery`, `AccountInfoQuery`, and `ContractInfoQuery` get all the data from consensus nodes again

## v2.38.0

### Added

-   `AccountBalanceQuery`, `AccountInfoQuery`, and `ContractInfoQuery` get part of the data from the Mirror Node REST API (HIP-367)
-   Fungible Token Metadata Field (HIP-646)
-   NFT Collection Token Metadata Field (HIP-765)

### Fixed

-   Raise an error if the transaction is not frozen while signing
-   Undeprecate `AccountBalance.TokenDecimals`, `AccountInfo.TokenRelationships`
-   Account alias for hollow account Mirror Node Queries

### Deprecated

-   `TokenRelationship.Symbol`, use `TokenInfo.Symbol` instead

## v2.37.0

### Added

-   METADATA key and possibility to update NFT metadata (HIP-657)
-   Fungible Token Metadata Field (HIP-646)
-   NFT Collection Token Metadata Field (HIP-765)
-   updated protobufs

## v2.36.0

### Added

-   Implemented HIP-844 (Add signerNonce field)

## v2.35.0

### Fixed

-   Implemented HIP-745 (Serialize transaction without freezing)

## v2.34.1

### Fixed

-   Fixed bug for fetching nodes from Client.

## v2.34.0

-   Refactored structures `executable`, `transaction` & `query` to define common methods in one place and reduce the code repetition

## v2.32.0

### Added

-   `PopulateAccount` to `AccountID`

## v2.31.0

### Added

-   `MarshalJSON` to `TransactionResponse`, `TransactionReceipt`, `TransactionRecord`
-   `IsApproved` boolean to `Transfer`

## v2.30.0

### Fixed

-   Node/Managed node concurrency issues
-   Serialization issues with `FileAppendTransaction`, `TokenCreateTransaction`, `TopicCreateTransaction`, `TopicMessageSubmitTransaction`

## v2.29.0

### Fixed

-   SDK Panics if node account id not in the network list
-   non-exhaustive switch statement errors for `AccountAllowanceApproveTransaction` and `AccountAllowanceDeleteTransaction`

## v2.28.0

### Added

-   `PopulateAccount` to `AccountID`
-   `PopulateContract` to `ContractID`

### Fixed

-   Data race in `TopicMessageQuery`

## v2.27.0

### Fixed

-   `ContractCreateFlow` to work with larger contract's bytecode

## v2.26.1

### Added

-   `ContractNonces` to `ContractFunctionResult` to support HIP-729

## v2.26.0

### Added

-   `AddInt*BigInt` functions to `ContractFunctionParameters` for sending `big.int`
-   `GetResult` function to `ContractFunctionResult` for parsing result to an interface
-   `GetBigInt` function to `ContractFunctionResult` for parsing result to a `big.int`

### Fixed

-   DER and PEM formats for private and public keys
-   Network concurrency issues
-   Some `ContractFunctionParameters` were sent/received as wrong data/data type

## v2.25.0

### Added

-   Added logging functionality to Client, Transaction and Query

## v2.24.4

### Added

-   Option to create a client from config file without mirror network
-   Finished LegderID implementation

## v2.24.3

### Added

-   Comments on the rest of the functions

## v2.24.2

### Added

-   Comments on `client` and all `query` and `transaction` types

## v2.24.1

### Fixed

-   `TransactionID.String()` will now return the TransactionID string without trimming the leading zeroes

## v2.24.0

### Added

-   Alias support in `AccountCreateTransaction`
-   `CreateAccountWithAliasExample`
-   `CreateAccountWithAliasAndReceiverSignatureRequiredExample`

## v2.23.0

### Added

-   `CustomFractionalFee.SetAssessmentMethod()`
-   `AccountAllowanceApproveTransaction.ApproveTokenNftAllowanceWithDelegatingSpender()`
-   `PrivateKeyFromStringECDSA()`
-   New mirror node endpoints, only 443 port is now supported `mainnet-public.mirrornode.hedera.com:443`,`testnet.mirrornode.hedera.com:443`, `previewnet.mirrornode.hedera.com:443`

### Fixed

-   Minimum query cost can now be less than 25 tinybars
-   `TransactionFromBytes()` now correctly sets transactionValidDuration

### Deprecated

-   `PrivateKeyFromStringECSDA()`

## v2.22.0

### Added

-   Support for HIP-583
-   Example for HIP-583 `account_create_token_transfer` which autocreates an account by sending HBAR to an Etherum Account Address
-   `Mnemonic.ToStandardEd25519PrivateKey` which uses the correct derivation path
-   `Mnemonic.ToStandardECDSAsecp256k1PrivateKey` which uses the correct derivation path

### Deprecated

-   `Mnemonic.ToPrivateKey()` was using incorrect derivation path
-   `PrivateKeyFromMnemonic()` was using incorrect derivation path

## v2.21.0

### Fixed

-   Panic in multiple node account id locking
-   Regenerating transaction id, when not expired
-   Signs more than once per node/transaction
-   Not retrying other nodes when there are multiple nodes are locked
-   Panic when locking multiple nodes
-   Panic when too many nodes are not healthy
-   INVALID_NODE_ACCOUNT error
-   Setting MaxAutomaticTokenAssociations on ContractUpdate even if not set

## v2.20.0

### Added

-   `IsZero()` and `Equals()` in `AccountID`
-   `GetSignedTransactionBodyBytes()` in `Transaction`

## v2.19.0

### Added

-   `SetDefaultMaxQueryPayment` and `SetDefaultMaxTransactionFee` in `Client`

### Fixed

-   Schedule network recursive call
-   Wrongly deprecated `SpenderID` field in `TokenNFTInfo`

## v2.18.0

### Added

-   `CustomFee.AllCollectorsAreExempt`

### Fixed

-   Addressbook Query no logner panics

## v2.17.7

### Added

-   Example for HIP-542
-   Example for HIP-573

### Fixed

-   10 seconds blocking on creating a `Client`

## v2.17.6

### Fixed

-   `NewAccountBallanceQuery` to correctly return account balance

## v2.17.5

### Added

-   Update documentation for `ContractFunction[Parameters|Result]` to show how to use large integer variants
-   Implement `TransactionResponse.ValidateStatus`

## v2.17.4

### Fixed

-   `*Transactions` now don't sign twice

## v2.17.3

### Added

-   `AccountCreateTransaction.[Set|Get]Alias[Key|EvmAddress]()`
-   `ContractCreateFlow.[Set|Get]MaxChunks()`
-   Support for automatically updating networks
-   `Client.[Set|Get]NetworkUpdatePeriod()`
-   `Client` constructor supports `_ScheduleNetworkUpdate` to disable auto updates
-   `task update` for manual address book updating

## v2.17.2

### Fixed

-   Removed deprecated flags for wrongfully deprecated `Client.[Set|Get]LedgerID`
-   Added `SetGrpcDeadline` to `EthereumTransaction` and `TokenUpdateTransaction`
-   Deprecated `LiveHash.Duration` use `LiveHash.LiveHashDuration`
-   Added missing `LiveHashQuery.[Set|Get]MaxRetry`
-   Added missing `TopicInfoQuery.SetPaymentTransactionID`

## v2.17.1

### Deprecated

-   `AccountBalance.[tokens|tokenDecimals]` use a mirror node query instead
-   `AccountInfo.tokenRelationships` use a mirror node query instead
-   `ContractInfo.tokenRelationships` use a mirror node query instead
-   `TokenNftInfo.SpenderID` replaced by `TokenNftInfo.AllowanceSpenderAccountID`

### Fixed

-   `Token[Update|Create]Transaction.KycKey`
-   `TokenCreateTransaction.FreezeDefaul` wasn't being set properly.
-   Requests should retry on `PLATFORM_NOT_ACTIVE`

## v2.17.0

### Added

-   `PrngThansaction`
-   `TransactionRecord.PrngBytes`
-   `TransactionRecord.PrngNumber`
-   `task` runner support

## Deprecated

-   `ContractFunctionResult.ContractStateChanges` with no replacement.

## v2.17.0-beta.1

### Added

-   `PrngThansaction`
-   `TransactionRecord.PrngBytes`
-   `TransactionRecord.PrngNumber`

## Deprecated

-   `ContractFunctionResult.ContractStateChanges` with no replacement.

## v2.16.1

### Added

-   `StakingInfo.PendingHbarReward`

## v2.16.0

### Added

-   `StakingInfo`
-   `AccountCreateTransaction.[Set|Get]StakedAccountID`
-   `AccountCreateTransaction.[Set|Get]StakedNodeID`
-   `AccountCreateTransaction.[Set|Get]DeclineStakingReward`
-   `AccountInfo.StakingInfo`
-   `AccountUpdateTransaction.[Set|Get]StakedAccountID`
-   `AccountUpdateTransaction.[Set|Get]StakedNodeID`
-   `AccountUpdateTransaction.[Set|Get]DeclineStakingReward`
-   `AccountUpdateTransaction.ClearStaked[AccountID|NodeID]`
-   `ContractCreateTransaction.[Set|Get]StakedNodeAccountID`
-   `ContractCreateTransaction.[Set|Get]StakedNodeID`
-   `ContractCreateTransaction.[Set|Get]DeclineStakingReward`
-   `ContractInfo.StakingInfo`
-   `ContractUpdateTransaction.[Set|Get]StakedNodeAccountID`
-   `ContractUpdateTransaction.[Set|Get]StakedNodeID`
-   `ContractUpdateTransaction.[Set|Get]DeclineStakingReward`
-   `ContractUpdateTransaction.ClearStaked[AccountID|NodeID]`
-   `TransactionRecord.PaidStakingRewards`
-   `ScheduleCreateTransaction.[Set|Get]ExpirationTime`
-   `ScheduleCreateTransaction.[Set|Get]WaitForExpiry`
-   Protobuf requests and responses will be logged, for `TRACE`, in hex.

### Fixed

-   `TopicMessageSubmitTransaction` empty `ChunkInfo` would always cause an error

## v2.16.0-beta.1

### Added

-   `StakingInfo`
-   `AccountCreateTransaction.[Set|Get]StakedNodeAccountID`
-   `AccountCreateTransaction.[Set|Get]StakedNodeID`
-   `AccountCreateTransaction.[Set|Get]DeclineStakingReward`
-   `AccountInfo.StakingInfo`
-   `AccountUpdateTransaction.[Set|Get]StakedNodeAccountID`
-   `AccountUpdateTransaction.[Set|Get]StakedNodeID`
-   `AccountUpdateTransaction.[Set|Get]DeclineStakingReward`
-   `ContractCreateTransaction.[Set|Get]StakedNodeAccountID`
-   `ContractCreateTransaction.[Set|Get]StakedNodeID`
-   `ContractCreateTransaction.[Set|Get]DeclineStakingReward`
-   `ContractInfo.StakingInfo`
-   `ContractUpdateTransaction.[Set|Get]StakedNodeAccountID`
-   `ContractUpdateTransaction.[Set|Get]StakedNodeID`
-   `ContractUpdateTransaction.[Set|Get]DeclineStakingReward`
-   `TransactionRecord.PaidStakingRewards`
-   `ScheduleCreateTransaction.[Set|Get]ExpirationTime`
-   `ScheduleCreateTransaction.[Set|Get]WaitForExpiry`

## v2.15.0

### Added

-   `EthereumFlow`
-   `EthereumTransactionData`

### Fixed

-   `Transaction.[From|To]Bytes` would ignore some variables
-   Fixed naming for `Ethereum.SetCallDataFileID()` and `Ethereum.SetMaxGasAllowanceHbar()` to be consistent with other sdks.

# v2.14.0

### Added

-   `ContractCreateTransaction.[Get|Set]MaxAutomaticTokenAssociations()`
-   `ContractCreateTransaction.[Get|Set]AutoRenewAccountId()`
-   `ContractCreateTransaction.[Get|Set]Bytecode()`
-   `ContractUpdateTransaction.[Get|Set]MaxAutomaticTokenAssociations()`
-   `ContractUpdateTransaction.[Get|Set|clear]AutoRenewAccountId()`
-   `ContractCreateFlow.[Get|Set]MaxAutomaticTokenAssociations()`
-   `ContractCreateFlow.[Get|Set]AutoRenewAccountId()`
-   `ContractInfo.AutoRenewAccountID`
-   `ContractDeleteTransaction.[Get|Set]PermanentRemoval`
-   `ContractCallQuery.[Get|Set]SenderID`
-   `ScheduleCreateTransaction.[Get|Set]ExpirationTime`
-   `ScheduleCreateTransaction.[Get|Set]WaitForExpiry`
-   `ScheduleInfo.WaitForExpiry`
-   `EthereumTransaction`
-   `TransactionRecord.EthereumHash`
-   `AccountInfo.EthereumNonce`
-   `AccountID.AliasEvmAddress`
-   `AccountID.AccountIDFromEvmAddress()`
-   `TransactionResponse.Get[Record|Receipt]Query`

## v2.14.0-beta.3

### Fixed

-   `FileUpdateTransaction` and `TopicMessageSubmitTransaction` duplicate transaction errors.
-   `*Transaction.ToBytes()` now properly chunked transactions.

## v2.14.0-beta.2

### Fixed

-   `*Query` payment signatures weren't getting updated after updating body with new random node.

# v2.14.0-beta.1

### Added

-   `ContractCreateTransaction.[Get|Set]MaxAutomaticTokenAssociations()`
-   `ContractCreateTransaction.[Get|Set]AutoRenewAccountId()`
-   `ContractCreateTransaction.[Get|Set]Bytecode()`
-   `ContractUpdateTransaction.[Get|Set]MaxAutomaticTokenAssociations()`
-   `ContractUpdateTransaction.[Get|Set|clear]AutoRenewAccountId()`
-   `ContractCreateFlow.[Get|Set]MaxAutomaticTokenAssociations()`
-   `ContractCreateFlow.[Get|Set]AutoRenewAccountId()`
-   `ContractInfo.AutoRenewAccountID`
-   `ContractDeleteTransaction.[Get|Set]PermanentRemoval`
-   `ContractCallQuery.[Get|Set]SenderID`
-   `ScheduleCreateTransaction.[Get|Set]ExpirationTime`
-   `ScheduleCreateTransaction.[Get|Set]WaitForExpiry`
-   `ScheduleInfo.WaitForExpiry`
-   `EthereumTransaction`
-   `TransactionRecord.EthereumHash`
-   `AccountInfo.EthereumNonce`
-   `AccountID.AliasEvmAddress`
-   `AccountID.AccountIDFromEvmAddress()`

## v2.13.4

### Fixed

-   `FileUpdateTransaction` and `TopicMessageSubmitTransaction` duplicate transaction errors.
-   `*Transaction.ToBytes()` now properly chunked transactions.

## v2.13.3

### Fixed

-   `*Query` payment signatures weren't getting updated after updating body with new random node.

## v2.13.2

### Added

-   `*Query.GetMaxQueryPayment()`
-   `*Query.GetQueryPayment()`

### Fixed

-   `*Query.GetPaymentTransactionID()` panic when not set.`
-   Removed unneeded parameter in `AccountDeleteTransaction.GetTransferAccountID()`.
-   `FileUpdateTransaction.GeFileMemo()` is now `FileUpdateTransaction.GetFileMemo()`
-   `TopicMessageSubmitTransaction` failing to send all messages, instead was getting duplicated transaction error.
-   `TopicMessageSubmitTransaction` would panic if no message was set.

## v2.13.1

### Added

-   `TokenNftAllowance.DelegatingSpender`
-   `AccountAllowanceApproveTransaction.AddAllTokenNftApprovalWithDelegatingSpender()`
-   `AccountAllowanceApproveTransaction.ApproveTokenNftAllowanceAllSerialsWithDelegatingSpender()`

## Deprecated

-   `AccountAllowanceAdjustTransaction` with no replacement.
-   `AccountAllowanceDeleteTransaction.DeleteAllTokenAllowances()` with no replacement.
-   `AccountAllowanceDeleteTransaction.DeleteAllHbarAllowances()` with no replacement.
-   `AccountInfo.[Hbar|Toke|Nft]Allowances`, with no replacement.
-   `TransactionRecord.[Hbar|Toke|Nft]Allowances`, with no replacement.

## v2.13.0

### Added

-   `AccountAllowanceDeleteTransaction`
-   `ContractFunctionResult.[gas|hbarAmount|contractFunctionParametersBytes]`
-   `AccountAllowanceExample`
-   `ScheduleTransferExample`

### Deprecated

-   `AccountAllowanceAdjustTransaction.revokeTokenNftAllowance()` with no replacement.
-   `AccountAllowanceApproveTransaction.AddHbarApproval()`, use `ApproveHbarAllowance()` instead.
-   `AccountAllowanceApproveTransaction.ApproveTokenApproval()`, use `GrantTokenNftAllowance()` instead.
-   `AccountAllowanceApproveTransaction.ApproveTokenNftApproval()`, use `ApproveTokenNftAllowance()` instead.

### Fixed

-   `*Transaction.GetTransactionID()` panic when not set.
-   `Transaction.Freeze()` now properly sets NodeAccountIDs
-   `*Query` payment transaction now properly contains the right NodeAccountIDs.

## v2.13.0-beta.1

### Added

-   `AccountAllowanceDeleteTransaction`
-   `ContractFunctionResult.[gas|hbarAmount|contractFunctionParametersBytes]`
-   `AccountAllowanceExample`
-   `ScheduleTransferExample`

### Deprecated

-   `AccountAllowanceAdjustTransaction.revokeTokenNftAllowance()` with no replacement.
-   `AccountAllowanceApproveTransaction.AddHbarApproval()`, use `ApproveHbarAllowance()` instead.
-   `AccountAllowanceApproveTransaction.ApproveTokenApproval()`, use `GrantTokenNftAllowance()` instead.
-   `AccountAllowanceApproveTransaction.ApproveTokenNftApproval()`, use `ApproveTokenNftAllowance()` instead.

## v2.12.0

### Added

-   `AccountInfoFlowVerify[Signature|Transaction]()`
-   `Client.[Set|Get]NodeMinReadmitPeriod()`
-   Support for using any node from the entire network upon execution
    if node account IDs have no been locked for the request.
-   Support for all integer widths for `ContractFunction[Result|Selector|Params]`

### Fixed

-   Ledger ID checksums
-   `TransactionFromBytes()` should validate all the transaction bodies are the same

### Changed

-   Network behavior to follow a more standard approach (remove the sorting we
    used to do).

## v2.12.0-beta.1

### Added

-   `AccountInfoFlowVerify[Signature|Transaction]()`
-   `Client.[Set|Get]NodeMinReadmitPeriod()`
-   Support for using any node from the entire network upon execution
    if node account IDs have no been locked for the request.
-   Support for all integer widths for `ContractFunction[Result|Selector|Params]`

### Fixed

-   Ledger ID checksums
-   `TransactionFromBytes()` should validate all the transaction bodies are the same

### Changed

-   Network behavior to follow a more standard approach (remove the sorting we
    used to do).

## v2.11.0

### Added

-   `ContractCreateFlow`
-   `Query.[Set|Get]PaymentTransactionID`
-   Verbose logging using zerolog
-   `*[Transaction|Query].[Set|Get]GrpcDeadline()`
-   `TransactionRecord.[hbar|Token|TokenNft]AllowanceAdjustments`
-   `TransferTransaction.AddApproved[Hbar|Token|Nft]Transfer()`
-   `AccountAllowanceApproveTransaction.Approve[Hbar|Token|TokenNft]Allowance()`
-   `AccountAllowanceAdjustTransaction.[Grant|Revoke][Hbar|Token|TokenNft]Allowance()`
-   `AccountAllowanceAdjustTransaction.[Grant|Revoke]TokenNftAllowanceAllSerials()`

### Fixed

-   `HbarAllowance.OwnerAccountID`, wasn't being set.
-   Min/max backoff for nodes should start at 8s to 60s
-   The current backoff for nodes should be used when sorting inside of network
    meaning nodes with a smaller current backoff will be prioritized
-   `TopicMessageQuery` start time should have a default

### Deprecated

-   `AccountUpdateTransaction.[Set|Get]AliasKey`

### Removed

-   `Account[Approve|Adjust]AllowanceTransaction.Add[Hbar|Token|TokenNft]AllowanceWithOwner()`

## v2.11.0-beta.1

### Added

-   `ContractCreateFlow`
-   `Account[Approve|Adjust]AllowanceTransaction.add[Hbar|Token|TokenNft]AllowanceWithOwner()`
-   `Query.[Set|Get]PaymentTransactionID`
-   Verbose logging using zerolog
-   `*[Transaction|Query].[Set|Get]GrpcDeadline()`

### Fixed

-   `HbarAllowance.OwnerAccountID`, wasn't being set.
-   Min/max backoff for nodes should start at 8s to 60s
-   The current backoff for nodes should be used when sorting inside of network
    meaning nodes with a smaller current backoff will be prioritized

### Deprecated

-   `AccountUpdateTransaction.[Set|Get]AliasKey`

## v2.10.0

### Added

-   `owner` field to `*Allowance`.
-   Added free `AddressBookQuery`.

### Fixed

-   Changed mirror node port to correct one, 443.
-   Occasional ECDSA invalid length error.
-   ContractIDFromString() now sets EvmAddress correctly to nil, when evm address is not detected

## v2.10.0-beta.1

### Added

-   `owner` field to `*Allowance`.
-   Added free `AddressBookQuery`.

### Fixed

-   Changed mirror node port to correct one, 443.

## v2.9.0

### Added

-   CREATE2 Solidity addresses can now be represented by a `ContractId` with `EvmAddress` set.
-   `ContractId.FromEvmAddress()`
-   `ContractFunctionResult.StateChanges`
-   `ContractFunctionResult.EvmAddress`
-   `ContractStateChange`
-   `StorageChange`
-   New response codes.
-   `ChunkedTransaction.[Set|Get]ChunkSize()`, and changed default chunk size for `FileAppendTransaction` to 2048.
-   `AccountAllowance[Adjust|Approve]Transaction`
-   `AccountInfo.[hbar|token|tokenNft]Allowances`
-   `[Hbar|Token|TokenNft]Allowance`
-   `[Hbar|Token|TokenNft]Allowance`
-   `TransferTransaction.set[Hbar|Token|TokenNft]TransferApproval()`

### Fixed

-   Requests not cycling though nodes.
-   Free queries not attempting to retry on different nodes.

### Deprecated

-   `ContractId.FromSolidityAddress()`, use `ContractId.FromEvmAddress()` instead.
-   `ContractFunctionResult.CreatedContractIDs`.

## v2.9.0-beta.2

### Added

-   CREATE2 Solidity addresses can now be represented by a `ContractId` with `EvmAddress` set.
-   `ContractId.FromEvmAddress()`
-   `ContractFunctionResult.StateChanges`
-   `ContractFunctionResult.EvmAddress`
-   `ContractStateChange`
-   `StorageChange`
-   New response codes.
-   `ChunkedTransaction.[Set|Get]ChunkSize()`, and changed default chunk size for `FileAppendTransaction` to 2048.
-   `AccountAllowance[Adjust|Approve]Transaction`
-   `AccountInfo.[hbar|token|tokenNft]Allowances`
-   `[Hbar|Token|TokenNft]Allowance`
-   `[Hbar|Token|TokenNft]Allowance`
-   `TransferTransaction.set[Hbar|Token|TokenNft]TransferApproval()`

### Fixed

-   Requests not cycling though nodes.
-   Free queries not attempting to retry on different nodes.

### Deprecated

-   `ContractId.FromSolidityAddress()`, use `ContractId.FromEvmAddress()` instead.
-   `ContractFunctionResult.CreatedContractIDs`.

## v2.9.0-beta.1

### Added

-   CREATE2 Solidity addresses can now be represented by a `ContractId` with `EvmAddress` set.
-   `ContractId.FromEvmAddress()`
-   `ContractFunctionResult.StateChanges`
-   `ContractFunctionResult.EvmAddress`
-   `ContractStateChange`
-   `StorageChange`
-   New response codes.
-   `ChunkedTransaction.[Set|Get]ChunkSize()`, and changed default chunk size for `FileAppendTransaction` to 2048.
-   `AccountAllowance[Adjust|Approve]Transaction`
-   `AccountInfo.[hbar|token|tokenNft]Allowances`
-   `[Hbar|Token|TokenNft]Allowance`
-   `[Hbar|Token|TokenNft]Allowance`
-   `TransferTransaction.set[Hbar|Token|TokenNft]TransferApproval()`

### Fixed

-   Requests not cycling though nodes.
-   Free queries not attempting to retry on different nodes.

### Deprecated

-   `ContractId.FromSolidityAddress()`, use `ContractId.FromEvmAddress()` instead.
-   `ContractFunctionResult.CreatedContractIDs`.

## v2.8.0

### Added

-   Support for regenerating transaction IDs on demand if a request
    responses with `TRANSACITON_EXPIRED`

## v2.8.0-beta.1

### Added

-   Support for regenerating transaction IDs on demand if a request
    responses with `TRANSACITON_EXPIRED`

## v2.7.0

### Added

-   `AccountId.AliasKey`, including `AccountId.[From]String()` support.
-   `[PublicKey|PrivateKey].ToAccountId()`.
-   `AliasKey` fields in `TransactionRecord` and `AccountInfo`.
-   `Nonce` field in `TransactionId`, including `TransactionId.[set|get]Nonce()`
-   `Children` fields in `TransactionRecord` and `TransactionReceipt`
-   `Duplicates` field in `TransactionReceipt`
-   `[TransactionReceiptQuery|TransactionRecordQuery].[Set|Get]IncludeChildren()`
-   `TransactionReceiptQuery.[Set|Get]IncludeDuplicates()`
-   New response codes.
-   Support for ECDSA SecP256K1 keys.
-   `PrivateKeyGenerate[ED25519|ECDSA]()`
-   `[Private|Public]KeyFrom[Bytes|String][DER|ED25519|ECDSA]()`
-   `[Private|Public]Key.[Bytes|String][Raw|DER]()`
-   `DelegateContractId`
-   `*Id.[from|to]SolidityAddress()`

### Deprecated

-   `PrivateKeyGenerate()`, use `PrivateKeyGenerate[ED25519|ECDSA]()` instead.

## v2.7.0-beta.1

### Added

-   `AccountId.AliasKey`, including `AccountId.[From]String()` support.
-   `[PublicKey|PrivateKey].ToAccountId()`.
-   `AliasKey` fields in `TransactionRecord` and `AccountInfo`.
-   `Nonce` field in `TransactionId`, including `TransactionId.[set|get]Nonce()`
-   `Children` fields in `TransactionRecord` and `TransactionReceipt`
-   `Duplicates` field in `TransactionReceipt`
-   `[TransactionReceiptQuery|TransactionRecordQuery].[Set|Get]IncludeChildren()`
-   `TransactionReceiptQuery.[Set|Get]IncludeDuplicates()`
-   New response codes.
-   Support for ECDSA SecP256K1 keys.
-   `PrivateKeyGenerate[ED25519|ECDSA]()`
-   `[Private|Public]KeyFrom[Bytes|String][DER|ED25519|ECDSA]()`
-   `[Private|Public]Key.[Bytes|String][Raw|DER]()`

### Deprecated

-   `PrivateKeyGenerate()`, use `PrivateKeyGenerate[ED25519|ECDSA]()` instead.

## v2.6.0

### Added

-   New smart contract response codes

### Deprecated

-   `ContractCallQuery.[Set|Get]MaxResultSize()`
-   `ContractUpdateTransaction.[Set|Get]ByteCodeFileID()`

## v2.6.0-beta.1

### Added

-   New smart contract response codes

### Deprecated

-   `ContractCallQuery.[Set|Get]MaxResultSize()`
-   `ContractUpdateTransaction.[Set|Get]ByteCodeFileID()`

## v2.5.1

### Fixed

-   `TransferTransaction.GetTokenTransfers()`
-   `TransferTransaction.AddTokenTransfer()`
-   Persistent error not being handled correctly
-   `TransactionReceiptQuery` should return even on a bad status codes.
    Only \*.GetReceipt()`should error on non`SUCCESS` status codes

## v2.5.1-beta.1

### Changed

-   Refactored and updated node account ID handling to err whenever a node account ID of 0.0.0 is being set

## v2.5.0-beta.1

### Deprecated

-   `ContractCallQuery.[Set|Get]MaxResultSize()`
-   `ContractUpdateTransaction.[Set|Get]ByteCodeFileID()`

## v2.5.0

### Fixed

-   `TransactionReceiptQuery` should fill out `TransactionReceipt` even when a bad `Status` is returned

## v2.4.1

### Fixed

-   `TransferTransaction` should serialize the transfers list deterministically

## v2.4.0

### Added

-   Support for toggling TLS for both mirror network and services network

## v2.3.0

### Added

-   `FreezeType`
-   `FreezeTransaction.[get|set]FreezeType()`

## v2.3.0-beta 1

### Added

-   Support for HIP-24 (token pausing)
    -   `TokenInfo.PauseKey`
    -   `TokenInfo.PauseStatus`
    -   `TokenCreateTransaction.PauseKey`
    -   `TokenUpdateTransaction.PauseKey`
    -   `TokenPauseTransaction`
    -   `TokenUnpauseTransaction`

## v2.2.0

### Added

-   Support for automatic token associations
    -   `TransactionRecord.AutomaticTokenAssociations`
    -   `AccountInfo.MaxAutomaticTokenAssociations`
    -   `AccountCreateTransaction.MaxAutomaticTokenAssociations`
    -   `AccountUpdateTransaction.MaxAutomaticTokenAssociations`
    -   `TokenRelationship.AutomaticAssociation`
    -   `TokenAssociation`
-   `Transaction*` helper methods - should make it easier to use the result of `TransactionFromBytes()`

### Fixed

-   TLS now properly confirms certificate hashes
-   `TokenUpdateTransaction.GetExpirationTime()` returns the correct time
-   Several `*.Get*()` methods required a parameter similiar to `*.Set*()`
    This has been changed completely instead of deprecated because we treated this as hard bug
-   Several `nil` dereference issues related to to/from protobuf conversions

## v2.2.0-beta.1

### Added

-   Support for automatic token associations
    -   `TransactionRecord.AutomaticTokenAssociations`
    -   `AccountInfo.MaxAutomaticTokenAssociations`
    -   `AccountCreateTransaction.MaxAutomaticTokenAssociations`
    -   `AccountUpdateTransaction.MaxAutomaticTokenAssociations`
    -   `TokenRelationship.AutomaticAssociation`
    -   `TokenAssociation`

## v2.1.16

### Added

-   Support for TLS
-   Setters which follow the builder pattern to `Custom*Fee`
-   `Client.[min|max]Backoff()` support

### Deprecated

-   `TokenNftInfoQuery.ByNftID()` - use `TokenNftInfoQuery.SetNftID()` instead
-   `TokenNftInfoQuery.[By|Set|Get]AccountId()` with no replacement
-   `TokenNftInfoQuery.[By|Set|Get]TokenId()` with no replacement
-   `TokenNftInfoQuery.[Set|Get]Start()` with no replacement
-   `TokenNftInfoQuery.[Set|Get]End()` with no replacement

## v2.1.15

### Fixed

-   `AssessedCustomFee.PayerAccountIDs` was misspelled

## v2.1.14

### Added

-   Support for `CustomRoyaltyFee`
-   Support for `AssessedCustomFee.payerAccountIds`

### Fixed

-   `nil` dereference issues within `*.validateNetworkIDs()`

## v2.1.13

### Added

-   Implement `Client.pingAll()`
-   Implement `Client.SetAutoChecksumValidation()` which validates all entity ID checksums on requests before executing

### Fixed

-   nil dereference errors when decoding invalid PEM files

## v2.1.12

### Added

-   Updated `Status` with new response codes
-   Support for `Hbar.[from|to]String()` to be reversible

## v2.1.11

### Removed

-   `*.AddCustomFee()` use `*.SetCustomFees()` instead

### Changes

-   Update `Status` with new codes

### Fixes

-   `PrivateKey.LegacyDerive()` should correctly handle indicies

## v2.1.11-beta.1

### Added

-   Support for NFTS
    -   Creating NFT tokens
    -   Minting NFTs
    -   Burning NFTs
    -   Transfering NFTs
    -   Wiping NFTs
    -   Query NFT information
-   Support for Custom Fees on tokens:
    -   Setting custom fees on a token
    -   Updating custom fees on an existing token

## v2.1.10

### Added

-   All requests should retry on gRPC error `INTERNAL` if the message contains `RST_STREAM`
-   `AccountBalance.Tokens` as a replacement for `AccountBalance.Token`
-   `AccountBalance.TokenDecimals`
-   All transactions will now `sign-on-demand` which should result in improved performance

### Fixed

-   `TopicMessageQuery` not calling `Unsubscribe` when a stream is cancelled
-   `TopicMessageQuery` should add 1 nanosecond to the `StartTime` of the last received message
-   `TopicMessageQuery` allocate space for entire chunked message ahead of time
    for retries
-   `TokenDeleteTransaction.SetTokenID()` incorrectly setting `tokenID` resulting in `GetTokenID()` always returning an empty `TokenID`
-   `TransferTransaction.GetTokenTransfers()` incorrectly setting an empty value

### Deprecated

-   `AccountBalance.Token` use `AccountBalance.Tokens` instead

## v2.1.9

### Fixed

-   `Client.SetMirroNetwork()` producing a nil pointer exception on next use of a mirror network
-   Mirror node TLS no longer producing nil pointer exception

## v2.1.8

### Added

-   Support TLS for mirror node connections.
-   Support for entity ID checksums which are validated whenever a request begins execution.
    This includes the IDs within the request, the account ID within the transaction ID, and
    query responses will contain entity IDs with a checksum for the network the query was executed on.

### Fixed

-   `TransactionTransaction.AddHbarTransfer()` incorrectly determine total transfer per account ID

## v2.1.7

### Fixed

-   `TopicMessageQuery.MaxBackoff` was not being used at all
-   `TopicMessageQuery.Limit` was being incorrectly update with full `TopicMessages` rather than per chunk
-   `TopicMessageQuery.StartTime` was not being updated each time a message was received
-   `TopicMessageQuery.CompletionHandler` was be called at incorrect times
-   Removed the use of locks and `sync.Map` within `TopicMessageQuery` as it is unncessary
-   Added default logging to `ErrorHandler` and `CompletionHandler`

## v2.1.6

-   Support for `MaxBackoff`, `MaxAttempts`, `RetryHandler`, and `CompletionHandler` in `TopicMessageQuery`
-   Default logging behavior to `TopicMessageQuery` if an error handler or completion handler was not set

### Fixed

-   Renamed `ScheduleInfo.Signers` -> `ScheduleInfo.Signatories`
-   `TopicMessageQuery` retry handling; this should retry on more gRPC errors
-   `TopicMessageQuery` max retry timeout; before this would could wait up to 4m with no feedback
-   `durationFromProtobuf()` incorrectly calculation duration
-   `*Token.GetAutoRenewPeriod()` and `*Token.GetExpirationTime()` nil dereference
-   `Hbar.As()` using multiplication instead of division, and should return a `float64`

### Added

-   Exposed `Hbar.Negated()`

## v2.1.5

###

-   Scheduled transaction support: `ScheduleCreateTransaction`, `ScheduleDeleteTransaction`, and `ScheduleSignTransaction`
-   Non-Constant Time Comparison of HMACs [NCC-E001154-006]
-   Decreased `CHUNK_SIZE` 4096->1024 and increased default max chunks 10->20

## v2.1.5-beta.5

### Fixed

-   Non-Constant Time Comparison of HMACs [NCC-E001154-006]
-   Decreased `CHUNK_SIZE` 4096->1024 and increased default max chunks 10->20
-   Renamed `ScheduleInfo.GetTransaction()` -> `ScheduleInfo.getScheduledTransaction()`

## v2.1.5-beta.4

### Fixed

-   `Transaction.Schedule()` should error when scheduling un-scheduable tranasctions

### Removed

-   `nonce` from `TransactionID`
-   `ScheduleTransactionBody` - should not be part of the public API

## v2.1.5-beta.3

### Fixed

-   `Transaction[Receipt|Record]Query` should not error for status `IDENTICAL_SCHEDULE_ALREADY_CREATED`
    because the other fields on the receipt are present with that status.
-   `ErrHederaReceiptStatus` should print `exception receipt status ...` instead of
    `exception precheck status ...`

## v2.1.5-beta.2

### Fixed

-   Executiong should retry on status `PLATFORM_TRANSACTION_NOT_CREATED`
-   Error handling throughout the SDK
    -   A precheck error shoudl be returned when the exceptional status is in the header
    -   A receipt error should be returned when the exceptional status is in the receipt
-   `TransactionRecordQuery` should retry on node precheck code `OK` only if we're not
    getting cost of query.
-   `Transaction[Receipt|Record]Query` should retry on both `RECEIPT_NOT_FOUND` and
    `RECORD_NOT_FOUND` status codes when node precheck code is `OK`

## v2.1.5-beta.1

### Fixed

-   Updated scheduled transaction to use new HAPI porotubfs

### Removed

-   `ScheduleCreateTransaction.AddScheduledSignature()`
-   `ScheduleCreateTransaction.GetScheduledSignatures()`
-   `ScheduleSignTransaction.addScheduledSignature()`
-   `ScheduleSignTransaction.GetScheduledSignatures()`

## v2.x

### Added

-   Support for scheduled transactions.
    -   `ScheduleCreateTransaction` - Create a new scheduled transaction
    -   `ScheduleSignTransaction` - Sign an existing scheduled transaction on the network
    -   `ScheduleDeleteTransaction` - Delete a scheduled transaction
    -   `ScheduleInfoQuery` - Query the info including `bodyBytes` of a scheduled transaction
    -   `ScheduleId`
-   Support for scheduled and nonce in `TransactionId`
    -   `TransactionIdWithNonce()` - Supports creating transaction ID with random bytes.
    -   `TransactionId.[Set|Get]Scheduled()` - Supports scheduled transaction IDs.
-   `TransactionIdWithValidStart()`

### Fixed

-   Updated protobufs [#120](https://github.com/hiero-ledger/hiero-sdk-go/issues/120)

### Deprecate

-   `NewTransactionId()` - Use `TransactionIdWithValidStart()` instead.

## v2.0.0

### Changes

-   All requests support getter methods as well as setters.
-   All requests support multiple node account IDs being set.
-   `TransactionFromBytes()` supports multiple node account IDs and existing
    signatures.
-   All requests support a max retry count using `SetMaxRetry()`
