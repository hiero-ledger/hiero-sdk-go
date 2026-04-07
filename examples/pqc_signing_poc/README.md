# Post-Quantum Cryptography (PQC) Signing POC

## Overview

This proof-of-concept demonstrates the use of **CRYSTALS-Dilithium** (Mode3 / ML-DSA-65), a NIST-standardized post-quantum digital signature algorithm, for signing Hedera transaction bytes. It compares Dilithium against the classical algorithms currently supported by the Hedera SDK: **Ed25519** and **ECDSA**.

## Why Post-Quantum Cryptography?

Quantum computers, once sufficiently powerful, will be able to break the elliptic curve and RSA-based cryptographic schemes that secure today's blockchain networks. Specifically:

- **Ed25519** and **ECDSA** rely on the hardness of the elliptic curve discrete logarithm problem (ECDLP)
- Shor's algorithm running on a cryptographically relevant quantum computer (CRQC) can solve ECDLP in polynomial time
- This would allow an attacker to derive private keys from public keys, forging signatures on any transaction

NIST has standardized several post-quantum algorithms as part of FIPS 203/204/205 (August 2024):
- **ML-KEM** (formerly CRYSTALS-Kyber) — key encapsulation
- **ML-DSA** (formerly CRYSTALS-Dilithium) — digital signatures ← **used in this POC**
- **SLH-DSA** (formerly SPHINCS+) — hash-based signatures

## What This POC Demonstrates

1. **Key generation** using Dilithium Mode3 (equivalent to ML-DSA-65, NIST security level 3)
2. **Signing** of simulated Hedera transaction body bytes
3. **Verification** of Dilithium signatures
4. **Size comparison** of keys and signatures across Ed25519, ECDSA, and Dilithium
5. **Performance benchmarks** comparing sign/verify operations
6. **SDK integration pattern** showing how Dilithium fits the `TransactionSigner` function signature

## What to Expect

### Key and Signature Sizes

| Algorithm            | Public Key | Private Key | Signature |
|---------------------|-----------|------------|-----------|
| Ed25519              | 32 bytes  | 64 bytes   | 64 bytes  |
| ECDSA P-256          | 33 bytes  | 32 bytes   | ~72 bytes |
| Dilithium3 (ML-DSA-65) | **1,952 bytes** | **4,000 bytes** | **3,293 bytes** |

Dilithium keys and signatures are **~50-60x larger** than Ed25519. This is the primary trade-off of lattice-based post-quantum cryptography.

### Performance

Dilithium signing and verification are surprisingly fast — often comparable to or faster than ECDSA, and within a small factor of Ed25519. Typical results on modern hardware:

| Operation     | Ed25519  | ECDSA P-256 | Dilithium3 |
|--------------|----------|-------------|------------|
| Key Generation | ~30 μs  | ~40 μs     | ~50 μs     |
| Signing       | ~50 μs  | ~30 μs     | ~100 μs    |
| Verification  | ~80 μs  | ~80 μs     | ~40 μs     |

> **Note**: Actual numbers vary by CPU. Run the benchmarks on your hardware for accurate results.

### Impact on Hedera Transactions

- **Bandwidth**: Transaction sizes will increase significantly (~3 KB per signature vs ~64 bytes)
- **Storage**: Each signed transaction on the ledger will consume more space
- **Protobuf changes**: The `SignaturePair` message in Hedera's protobuf schema currently expects Ed25519/ECDSA-sized signatures; new fields or oneof variants would be needed
- **Key encoding**: DER/PEM encoding standards for Dilithium keys are still being finalized (OIDs assigned by NIST in FIPS 204)
- **Gossip/consensus**: Larger signatures may affect network throughput during consensus

## How to Run

### Run the demo

```bash
go run ./examples/pqc_signing_poc/
```

### Run the benchmarks

```bash
go test ./examples/pqc_signing_poc/ -bench=. -benchmem
```

For more stable results with longer benchmark duration:

```bash
go test ./examples/pqc_signing_poc/ -bench=. -benchmem -benchtime=3s
```

### Compare specific operations

```bash
# Only signing benchmarks
go test ./examples/pqc_signing_poc/ -bench=BenchmarkSign -benchmem

# Only verification benchmarks
go test ./examples/pqc_signing_poc/ -bench=BenchmarkVerify -benchmem

# Only key generation benchmarks
go test ./examples/pqc_signing_poc/ -bench=BenchmarkKeyGen -benchmem
```

## Limitations

- **Not network-compatible**: The Hedera network does not currently accept Dilithium signatures. This POC only demonstrates the cryptographic operations locally.
- **No protobuf integration**: Dilithium keys/signatures are not yet encoded in Hedera's protobuf `Key` or `SignaturePair` messages.
- **Single algorithm**: Only Dilithium Mode3 is shown. A production implementation might also evaluate SLH-DSA (SPHINCS+) for stateless hash-based signatures or hybrid classical+PQC schemes.
- **No hybrid signing**: A migration strategy would likely involve hybrid signatures (Ed25519 + Dilithium) during a transition period.

## Next Steps for Production

1. **Protobuf schema changes**: Add `dilithium` or `ml_dsa_65` variants to the `Key` and `SignaturePair` oneofs
2. **SDK key types**: Implement `_DilithiumPrivateKey` / `_DilithiumPublicKey` following the existing Ed25519/ECDSA patterns
3. **Hybrid signatures**: Support dual-signing (classical + PQC) for backwards compatibility
4. **Key derivation**: Investigate BIP32-like HD wallet derivation for lattice-based keys
5. **Network consensus**: Evaluate impact of larger signatures on gossip protocol throughput
6. **HIP proposal**: Draft a Hedera Improvement Proposal for PQC support

For a detailed breakdown of all ecosystem changes required, see [ECOSYSTEM_CHANGES.md](ECOSYSTEM_CHANGES.md).

## Dependencies

- [`github.com/cloudflare/circl`](https://github.com/cloudflare/circl) — Cloudflare's cryptographic library implementing NIST post-quantum standards
  - `circl/sign/dilithium/mode3` — CRYSTALS-Dilithium Mode3 (equivalent to ML-DSA-65)

## References

- [NIST FIPS 204 — ML-DSA (Dilithium)](https://csrc.nist.gov/pubs/fips/204/final)
- [CRYSTALS-Dilithium specification](https://pq-crystals.org/dilithium/)
- [Cloudflare circl library](https://github.com/cloudflare/circl)
- [Hedera SDK signing documentation](https://docs.hedera.com/hedera)
