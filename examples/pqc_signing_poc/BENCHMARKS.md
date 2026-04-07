# PQC Benchmark Results

## Test Environment

- **CPU**: Apple M2 Pro (ARM64)
- **OS**: macOS (Darwin 24.5.0)
- **Go**: 1.24+
- **Library**: cloudflare/circl v1.6.3
- **Benchmark duration**: 2s per operation (`-benchtime=2s`)

## Results

### Key Generation

| Algorithm            |   Time (μs) | Allocs/op | Bytes/op |
|---------------------|-------------|-----------|----------|
| Ed25519              |       17.3  |         3 |      128 |
| ECDSA P-256          |       11.5  |        16 |      984 |
| Dilithium3 (ML-DSA-65) |  **147.4** |         3 |   81,952 |

Dilithium key generation is ~8.5x slower than Ed25519 and ~12.8x slower than ECDSA. At 147 μs this is negligible in practice — a Hedera transaction's network round-trip alone is 3-7 **seconds**.

### Signing

| Algorithm            |   Time (μs) | Allocs/op | Bytes/op |
|---------------------|-------------|-----------|----------|
| Ed25519              |       22.4  |         1 |       64 |
| ECDSA P-256          |       22.6  |        67 |    6,336 |
| Dilithium3 (ML-DSA-65) |  **235.5** |         1 |      448 |

Dilithium signing is ~10.5x slower than Ed25519. Still under 250 μs — well within acceptable latency for any transaction workflow.

### Verification

| Algorithm            |   Time (μs) | Allocs/op | Bytes/op |
|---------------------|-------------|-----------|----------|
| Ed25519              |       48.9  |         0 |        0 |
| ECDSA P-256          |       53.4  |        10 |      576 |
| Dilithium3 (ML-DSA-65) |   **43.8** |         1 |      448 |

**Dilithium verification is the fastest of all three algorithms.** This is significant because verification happens far more often than signing — every node in the network verifies every transaction signature.

### Memory Efficiency

| Algorithm            | KeyGen (B/op) | Sign (B/op) | Verify (B/op) |
|---------------------|--------------|-------------|---------------|
| Ed25519              |          128 |          64 |             0 |
| ECDSA P-256          |          984 |       6,336 |           576 |
| Dilithium3 (ML-DSA-65) |   81,952 |         448 |           448 |

Dilithium key generation allocates ~82 KB (due to the 4 KB private key + 1.9 KB public key + internal matrix expansion), but signing and verification are lean — fewer allocations than ECDSA.

## Size Comparison

| Algorithm            | Public Key | Private Key | Signature |
|---------------------|-----------|------------|-----------|
| Ed25519              |   32 B    |     64 B   |    64 B   |
| ECDSA P-256          |   33 B    |     32 B   |   ~70 B   |
| Dilithium3 (ML-DSA-65) | **1,952 B** | **4,000 B** | **3,293 B** |

- Public key: **61x** larger than Ed25519
- Signature: **51x** larger than Ed25519
- A single signed transaction grows by ~3.2 KB per Dilithium signature

## Throughput Estimates

Based on benchmarks, single-threaded throughput:

| Algorithm            | Signs/sec | Verifications/sec |
|---------------------|-----------|-------------------|
| Ed25519              |  ~44,700  |         ~20,400   |
| ECDSA P-256          |  ~44,300  |         ~18,700   |
| Dilithium3 (ML-DSA-65) | ~4,250 |       **~22,800** |

Dilithium can verify over 22,000 signatures per second per core — higher than both classical algorithms.

## How to Reproduce

```bash
# Full benchmark suite
go test ./examples/pqc_signing_poc/ -bench=. -benchmem -benchtime=2s

# Signing only
go test ./examples/pqc_signing_poc/ -bench=BenchmarkSign -benchmem -benchtime=3s

# Compare with CPU profile
go test ./examples/pqc_signing_poc/ -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## Key Takeaway

**Performance is not a blocker for PQC adoption.** The bottleneck will be bandwidth and storage (51-61x larger keys/signatures), not computation time. Dilithium's verification speed is actually an advantage for network nodes that must verify every transaction.
