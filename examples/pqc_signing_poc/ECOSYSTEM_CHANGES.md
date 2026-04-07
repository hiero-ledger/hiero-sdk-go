# Hedera Ecosystem Changes Required for PQC

Adopting PQC is not just an SDK change — it touches nearly every layer of the Hedera stack. Below is a concrete breakdown of what needs to change for production deployment.

## 1. Protobuf Schema (Blocking)

The `Key` and `SignaturePair` messages in `basic_types.proto` only have variants for Ed25519, ECDSA, and RSA. New `oneof` fields are required:

```protobuf
message Key {
  oneof key {
    // existing fields...
    bytes dilithium_mode3 = 10;   // ML-DSA-65 (field number TBD)
  }
}

message SignaturePair {
  oneof signature {
    // existing fields...
    bytes dilithium_mode3 = 8;    // ML-DSA-65 (field number TBD)
  }
}
```

Field numbers must be agreed across all implementations. The current POC works around this by stuffing Dilithium bytes into the Ed25519 field — **not valid for network use**.

## 2. Consensus Nodes (hedera-services)

The Java-based consensus nodes must be updated to:

- **Parse** the new protobuf key/signature variants
- **Verify** Dilithium signatures during transaction pre-check and consensus (note: Dilithium verification is actually **faster** than Ed25519 — see benchmarks)
- **Store** much larger keys and signatures in state (public key: 1,952 B vs 32 B — 61x)
- **Handle bandwidth**: a block of 1,000 transactions grows from ~64 KB to ~3.3 MB in signature data alone
- **Update fee schedules** to account for larger transaction sizes and storage

## 3. All 7 SDKs (Go, Java, JS/TS, Swift, Rust, C++, Python)

Each SDK needs:

- New key types wrapping a Dilithium (ML-DSA-65) implementation
- Key generation, serialization, and deserialization (raw bytes + DER with FIPS 204 OIDs)
- Integration into existing Sign / Verify / SignTransaction / VerifyTransaction flows
- **Hybrid signing** support (classical + PQC dual signatures for migration period)
- BIP32-like HD wallet derivation — **not yet standardized** for lattice-based keys; the seed can stay 32 bytes but derivation paths need a spec

## 4. Mirror Node & REST APIs

- Parse, index, and serve the new key/signature types
- Return PQC keys in account info, transaction records, and state proofs
- Handle significantly larger data volumes in state proofs (~33 KB for 10 node sigs vs ~640 B with Ed25519)
- REST API key representation: hex-encoding a 1,952-byte key = a 3,904-character string
- gRPC message size limits may need adjustment

## 5. Wallet & HSM Ecosystem

- **Hardware wallets** (Ledger, Trezor) — limited memory may not support 4 KB private keys or 3.3 KB signatures
- **Cloud KMS** (AWS, GCP, Azure) — no PQC signing support as of 2025
- **Browser wallets** (HashPack, Blade, etc.) — need JS/WASM Dilithium implementations
- **Enterprise HSMs** — most do not yet support ML-DSA natively; FIPS 140-3 validated PQC modules are not widely available

## 6. Governance & Migration

- A **HIP (Hedera Improvement Proposal)** defining the PQC rollout strategy
- **Hybrid-first migration**: Ed25519 + Dilithium dual signatures during transition, allowing accounts to opt-in before PQC-only enforcement
- **Protocol versioning**: clear rules for when PQC-only signatures are accepted
- **Rollback plan**: ability to strip PQC signatures if Dilithium is later broken, without invalidating transactions
- Network upgrade vote by council nodes

## Recommended Phased Approach

| Phase | Scope | Goal |
|-------|-------|------|
| **0 — POC** (this) | SDK-level crypto only | Validate algorithm, measure sizes/perf |
| **1 — HIP & Protobuf** | Draft HIP, agree on proto field numbers | Cross-SDK alignment |
| **2 — SDK Support** | All 7 SDKs add Dilithium key types | Client-side readiness |
| **3 — Hybrid on Network** | Nodes accept hybrid (classical + PQC) sigs | Backwards-compatible PQC |
| **4 — PQC-Only Option** | Accounts can use Dilithium-only keys | Full quantum resistance |
| **5 — Ecosystem Maturity** | Wallets, HSMs, KMS, DApps catch up | End-to-end PQC |
