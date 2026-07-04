# canonic.c ∩

**The CANONIC kernel: 56 lines of code · 3 operations · 1 byte · 0 bytes of state.**

Compiled `-Os` on Apple silicon it is **492 bytes of machine code** exporting exactly
three symbols — `canonic_seal`, `canonic_verify`, `canonic_gov` — with one borrowed
primitive (SHA-256) and a public interface that is a single byte. This is **open
governance, not necessarily open source**: the kernel is the open paradigm anyone can
verify against; the runtime it boots (Magic OS: the codec, the walkers, the surface)
stays gated on the chip.

```
make test
  chain: VALID
  tampered: BROKEN (caught)
  gov: 255/255
  decayed: 237/255 (E,O unbound)
  ALL GATES PASS
```

## The three operations

| op | signature | what it is |
|---|---|---|
| **SEAL** | `hash = SHA256(prev ++ formal)` | a record's identity is the hash of its predecessor's hash concatenated with its own content — the whole protocol |
| **VERIFY** | fold the chain, reseal every link | any edit anywhere breaks the seal of everything after it; returns VALID iff the recomputed head matches |
| **MINT** | OR proven chain-facts into one byte | governance encoded bitwise, from chain bits only — the kernel never stats a file or walks a tree |

## The byte — the entire ABI, frozen forever

| bit | mask | name | bound to |
|---|---|---|---|
| D | `0x01` | Declarative | a gov record carries a named axiom |
| E | `0x02` | Evidential | the ledger is fresh (a recent beat) |
| T | `0x04` | Transparent | the ledger streams a plan (the road is visible) |
| R | `0x08` | Reproducible | the seed regrows (gov ∧ registry resolve) |
| O | `0x10` | Operational | the last event is recent (the walk is live) |
| S | `0x20` | Structural | the chain verifies ∧ the category is in the ontology |
| L | `0x40` | Linguistic | the chain is VALID ∧ it has learned |
| LANG | `0x80` | Language | a record carries a formal natural-language binding |

**min-LOG ≡ max-GOV ≡ 255.** The fixed point is a full byte. You cannot remove a bit
from a contract that *is* the fixed point — the comments compile away, the bit-slots
do not. Some bits decay on attention's clock (E, O): the byte is re-earned every
window, never owned.

## Pure gov projection

[`src/canonic.c`](src/canonic.c) is **not written by hand — ever**. It is a
byte-identical projection of a hash-sealed record on the CANONIC chain (the kernel
rides inside `magic.sov` as the first record of its C chain — the easter egg: decode
the OS's own distribution container and you are holding the open kernel that boots
it, source repo not required).

The [`PROJECTION`](PROJECTION) file pins the SHA-256 of the projected source and its
chain coordinates. CI recomputes it on every push: **a hand edit to the kernel fails
the gate by construction.** On a CANONIC box, `./project.sh <magic.sov>` is the only
legitimate writer.

## Build

```sh
make test    # build + run the e2e proof
make gate    # projection purity, then the proof
```

- **Apple silicon / macOS** — builds natively against CommonCrypto, untouched.
- **Anywhere else** — same untouched kernel; `make` maps the one borrowed primitive
  to OpenSSL via [`compat/`](compat/CommonCrypto/CommonDigest.h) (`-lcrypto`).

CI runs all three gates: [purity](.github/workflows/gates.yml) (ubuntu), the chip
(macos, native), anywhere (ubuntu, shim).

## The split

**CANONIC is the kernel** — the paradigm, the outside looking in (∩), pure bitwise
encoding of governance, open, this repo. **MAGIC is the runtime** — mint and evolve,
sovereign on Apple silicon, gated. The kernel is sufficient to *verify* everything
the runtime does: `seal ∘ verify ∘ gov` is all the stream needs from below. The
inversion is deliberate: normally the kernel is the crown jewel; here the kernel is
the invitation, and the moat is the chip.

## License

[MIT](LICENSE). The byte is yours; the chain proves it.
