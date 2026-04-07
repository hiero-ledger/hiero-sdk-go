# PQC Risk Analysis for Hedera

Performance and key/signature size are the obvious trade-offs, but they are arguably the **least concerning** risks. Below is a broader risk assessment for adopting post-quantum cryptography in the Hedera ecosystem.

## 1. Cryptanalytic Maturity

**Risk: The algorithms are young.**

CRYSTALS-Dilithium was standardized by NIST in August 2024 (FIPS 204). While it survived multiple rounds of public cryptanalysis since 2017, lattice-based cryptography has not endured the **40+ years** of scrutiny that RSA and elliptic curves have.

- In 2022, a NIST PQC finalist (SIKE) was [completely broken](https://eprint.iacr.org/2022/975) after years of analysis
- Lattice problems (Module-LWE, Module-SIS) are believed hard, but "believed hard" is not "proven hard"
- New algebraic or quantum attacks on structured lattices remain an active research area

**Mitigation**: Use hybrid signatures (classical + PQC) during a transition period. If Dilithium is broken, the classical signature still protects.

## 2. "Harvest Now, Decrypt Later" Threat

**Risk: Data recorded on-chain today is exposed to future quantum attacks.**

Adversaries can record blockchain transactions now and attempt to extract private keys from public keys once quantum computers are available. This is especially relevant for:

- Long-lived accounts with publicly visible keys (e.g., treasury accounts, governance keys)
- Immutable on-chain data that reveals public keys
- Accounts that must remain secure for 10+ years

**Mitigation**: This is actually the strongest argument **for** early PQC adoption, even before quantum computers exist. Accounts with long-term security requirements should migrate first.

## 3. Implementation & Side-Channel Attacks

**Risk: Correct algorithms can be defeated by flawed implementations.**

Lattice-based schemes introduce new side-channel attack surfaces:

- **Timing attacks**: Rejection sampling in Dilithium signing can leak information if not constant-time
- **Power/EM analysis**: Particularly relevant for hardware wallets and HSMs signing with Dilithium
- **Fault injection**: Inducing faults during lattice operations can reveal secret key components
- **Nonce reuse**: Unlike Ed25519 (which derives nonces deterministically from the message), some Dilithium modes use randomized nonces — a faulty RNG during signing can be catastrophic

The `circl` library from Cloudflare is well-maintained, but it has not been through the same level of FIPS certification and hardware-level audit as OpenSSL's Ed25519/ECDSA.

**Mitigation**: Use audited, FIPS-validated implementations when available. Prefer deterministic signing modes. Test on target hardware platforms.

## 4. Ecosystem & Interoperability

**Risk: The broader ecosystem is not PQC-ready.**

Hedera does not exist in isolation. PQC adoption requires alignment across:

- **All 7 Hedera SDKs** (Go, Java, JS/TS, Swift, Rust, C++, Python) — each must add Dilithium support
- **Hardware wallets** (Ledger, Trezor) — limited memory/storage may not support 4 KB private keys
- **HSMs** — most enterprise HSMs do not yet support ML-DSA natively
- **Key management systems** (Vault, KMS, cloud providers) — AWS KMS, GCP KMS, Azure Key Vault have limited or no PQC support as of 2025
- **TLS/gRPC** — client-to-node communication; PQC TLS (using ML-KEM for key exchange) is still being rolled out
- **Wallets and DApps** — every wallet, explorer, and application that handles keys needs updating
- **Certificate authorities** — if Hedera uses TLS certificates, CA infrastructure must support PQC

A partial rollout where only some participants use PQC creates a fragmented security model.

**Mitigation**: Phased rollout. Start with SDK support and optional PQC key types. Network-level enforcement comes later.

## 5. Network & Consensus Impact

**Risk: Larger signatures affect throughput, not just storage.**

Hedera's hashgraph consensus propagates signed events through gossip. Larger signatures mean:

- **Gossip bandwidth**: Each event carries signatures; 51x larger signatures increase gossip traffic substantially
- **State proofs**: Hedera state proofs include signatures from 1/3+ of stake — if each is 3.3 KB instead of 64 bytes, state proofs grow by tens of KB
- **Transaction throughput**: Larger transactions mean fewer transactions per second at the same bandwidth
- **Mirror node storage**: Historical transaction records grow significantly
- **gRPC message limits**: Default protobuf/gRPC message size limits may need adjustment

**Rough estimate**: A Hedera state proof with 10 node signatures would grow from ~640 bytes of signatures to ~33 KB. A block of 1,000 transactions would carry ~3.3 MB of Dilithium signatures vs ~64 KB of Ed25519 signatures.

**Mitigation**: Consider signature aggregation research for lattice-based schemes (still an open research problem). Evaluate SLH-DSA (SPHINCS+) as an alternative with different size/speed trade-offs. Optimize protobuf encoding.

## 6. Key Management Complexity

**Risk: Larger keys create operational challenges.**

- **Backup and recovery**: A Dilithium private key is 4 KB — too large for QR codes, paper wallets, or manual transcription
- **Mnemonic derivation**: BIP39/BIP32 HD wallet derivation for Dilithium keys is not standardized; there is no equivalent of "12 words → private key" yet
- **Multi-signature**: Threshold signature schemes (TSS) for lattice-based cryptography are still experimental research
- **Key rotation**: Migrating accounts from Ed25519 to Dilithium keys requires careful coordination to avoid losing access

**Mitigation**: Keep seed-based derivation (the seed can remain 32 bytes even if the expanded key is 4 KB). Invest in UX for key migration workflows.

## 7. Regulatory & Compliance

**Risk: Regulatory requirements are evolving.**

- NIST has standardized ML-DSA, but **FIPS-validated implementations** are not yet widely available
- Some regulated industries (finance, healthcare) require FIPS 140-3 validated crypto modules — these won't be available for PQC immediately
- The EU's Cyber Resilience Act and other frameworks may mandate PQC readiness on a timeline that doesn't align with Hedera's roadmap
- Government agencies (NSA CNSA 2.0) are mandating PQC adoption by 2035 — this creates pressure but also a clear timeline

**Mitigation**: Track FIPS 140-3 validation timelines for PQC modules. Plan for compliance requirements in advance.

## 8. Hybrid Signature Complexity

**Risk: The transition period is the most dangerous phase.**

A realistic migration involves dual/hybrid signatures (Ed25519 + Dilithium) for backward compatibility. This introduces:

- **Double the signing cost** (though still under 300 μs combined)
- **Increased transaction size** (64 + 3,293 = 3,357 bytes per signature pair)
- **Complex verification logic**: Nodes must verify both signatures and handle mixed-key accounts
- **Protocol versioning**: Need clear rules for when PQC-only signatures are accepted
- **Rollback risk**: If Dilithium is later broken, need ability to strip PQC signatures without invalidating transactions

**Mitigation**: Design the hybrid scheme carefully upfront. Define clear protocol version gates.

## 9. Protobuf & Wire Format

**Risk: Hedera's protobuf schema assumes small keys and signatures.**

The current `Key` protobuf message uses a `oneof` with fields like `ed25519` (32 bytes) and `ECDSASecp256k1` (33 bytes). Adding Dilithium requires:

- New `oneof` variant for ML-DSA-65 keys (~1,952 bytes)
- Updated `SignaturePair` to handle 3,293-byte signatures
- Possible changes to `SignatureMap` encoding
- Mirror node API changes to return/accept PQC key types
- REST API changes for key representation (hex encoding a 1,952-byte key = 3,904 character string)

**Mitigation**: Plan protobuf schema evolution early. Use `bytes` fields with algorithm identifiers rather than fixed-size fields.

## 10. Algorithm Agility

**Risk: Locking into one PQC algorithm is itself a risk.**

If Dilithium is later found to have a weakness, Hedera needs the ability to migrate to another PQC scheme (e.g., SLH-DSA, Falcon, or a future algorithm). This requires:

- **Algorithm-agile key types**: Don't hardcode Dilithium; design for pluggable PQC algorithms
- **Key rotation mechanisms**: Accounts must be able to update their key type without losing identity
- **Version negotiation**: Nodes must agree on which algorithms are acceptable at any given protocol version

**Mitigation**: Design the SDK and protocol around algorithm identifiers (OIDs), not specific algorithm implementations.

## Risk Summary Matrix

| Risk | Likelihood | Impact | Urgency |
|------|-----------|--------|---------|
| Cryptanalytic break | Low | Critical | Low (use hybrid) |
| Harvest-now-decrypt-later | High | High | **High** (motivates early action) |
| Side-channel attacks | Medium | High | Medium |
| Ecosystem readiness | High | Medium | Medium |
| Network throughput | Medium | High | Medium |
| Key management UX | High | Medium | Low |
| Regulatory pressure | Medium | Medium | Low (2035 deadline) |
| Hybrid complexity | High | Medium | Medium |
| Protobuf changes | Certain | Medium | Medium |
| Algorithm agility | Medium | High | Medium |

## Bottom Line

**Speed is not the problem** — Dilithium is fast enough. The real challenges are:

1. **Size** — 51-61x larger keys/signatures impact the entire stack from gossip to storage
2. **Ecosystem** — wallets, HSMs, KMS, other SDKs, and DApps all need PQC support
3. **Migration** — the hybrid transition period introduces complexity and risk
4. **Maturity** — lattice-based cryptography is young; algorithm agility is essential
5. **Harvest-now-decrypt-later** — the strongest argument for acting sooner rather than later

The recommended approach is a **phased, hybrid-first migration** that adds PQC as an optional key type alongside existing Ed25519/ECDSA, with a clear roadmap toward PQC-only once the ecosystem matures.
