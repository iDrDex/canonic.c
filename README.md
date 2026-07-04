# canonic.c ∩ — sov code, not source code

**This branch carries no source tree. The artifact is [`canonic.sov`](canonic.sov) —
~8 KB that regrows, builds, and proves the entire project.**

`canonic.sov` is a tar.gz holding a manifest and one hash chain in which **every file
of the project is a sealed record**, chained by the kernel's own law —
`SHA256(prev ++ content)`, genesis all-zeros — with the 56-line CANONIC kernel itself
as record zero. File trees are *projections* of this sov, regrown on demand:

```sh
./regrow.sh --into grown     # sov chain: VALID · projection: REGROWN (12 files)
cd grown && make gate        # purity OK · VALID · tamper caught · 255/255 · ALL GATES PASS
```

CI does exactly that on every push, on ubuntu and macos: **8 KB in, a fully gated
build out.** The sov is sufficient; everything else is hospitality.

## The branches — every projection is a branch

| branch | what it is |
|---|---|
| **`main`** | the trunk: the sov + the regrow gate. Sov code. |
| [**`projection/src`**](../../tree/projection/src) | the C-tree lens — the sov regrown into browsable files, CI-gated **EXACT** against the sov (any drift between tree and sov fails the build) |

Change flows **sov → branch**, never branch → sov. A future face (a Go client, a TS
verifier) lands as another `projection/*` branch regrown from its own chain.

## The kernel (record zero)

**56 lines of code · 3 operations · 1 byte · 0 bytes of state — 492 bytes of ARM64.**
`canonic_seal` (hash = SHA256(prev ++ formal)) · `canonic_verify` (fold the chain,
reseal every link) · `canonic_gov` (OR proven chain-facts into the one-byte ABI:
D E T R O S L LANG, frozen — min-LOG ≡ max-GOV ≡ 255; you cannot remove a bit from a
contract that *is* the fixed point). Browse it on
[`projection/src`](../../tree/projection/src) or regrow it yourself from the sov.

**Open gov, not necessarily open source**: the kernel is the open paradigm anyone can
verify against; the runtime it boots (Magic OS — mint and evolve, sovereign on Apple
silicon) stays gated on the chip. The kernel is the invitation; the moat is the chip.

## License

[MIT](LICENSE). The byte is yours; the chain proves it.
