# CANONIC — THE BUILD TREE   (normalized · maximally declarative · the build IS the verbs)
#   the verbs tell the story; the CANONICS are the composed nouns; one noun per verb (the 1:1).

CRYPTO  = attest( ## )                            # TRUST  · no signature, no fold      [alg=p256-se]
LAMBDA  = nano( CRYPTO ∘ FORMAL )                 # WRITE  · the one chain writer       [writers=1 must be 1]
KILN    = encode( bin ∪ src ∪ SOV ∪ libexec )     # BUILD  · the codec seal             [root=0fb0006e… leaves=73]
INTEL   = coal( LAMBDA* )                         # BUILD  · the derived brain          [gap=2]
MAGIC   = laude( INTEL )   ≡   CODEX = sov( LAMBDA* )   # SPEAK · one interface, two directions [pair=ok]
CHAT    = drain( LAMBDA* )                        # SHIP   · federation, chain-down     [fleet=behind]
CODON   = qed( CRYPTO LAMBDA KILN INTEL MAGIC CODEX CHAT )  # PROVE · the byte          [byte=0xDF CS·NRXEP]

CANONIC = fix( CODON = 0xFF )                     # the UPPER build: lower builds UPPER

## BIN — the 8-verb kernel (bin ↔ CANONICS 1:1, one verb per qed bit)

| bit | verb | bin | atom (chain `#alias BIN:`) | byte echo |
|-----|------|-----|------------------------------|-----------|
| b7·C | FOLD | `nano` | LAMBDA | LAMBDA |
| b6·S | SIGN | `attest` | CRYPTO | CRYPTO |
| b5·B | BEAT | `coal` | INTEL | INTEL |
| b4·N | NAME | `libexec/nomen` | LANGUAGE | LANGUAGE |
| b3·R | RECALL | `laude` | MAGIC | MAGIC |
| b2·X | FLEET | `drain` | CHAT | CHAT |
| b1·E | BIND | `encode` | KILN | KILN |
| b0·P | RESOLVE | `sov` | CODEX | CODEX |
| byte | QED | `qed` | CODON | CODON |

## STRATA — the 5 build commands (lower builds UPPER; qed = the byte over the tower)

| # | STRATUM | verbs (lower) | measured invariant |
|---|---------|---------------|--------------------|
| 0 | TRUST | `attest` | SE hardware: every write depends on it, it depends on nothing |
| 1 | WRITE | `nano` | sole chain writer — files appending to the ledger: **1** (must be 1, the engine) |
| 2 | BUILD | `encode` · `coal` | encode is the leaf (kernel refs in bytes: **0**, must be 0); coal derives, never writes |
| 3 | SPEAK | `laude` ≡ `sov` | LAUDE IS SOV — mutual exec measured: **one-interface** (meaning→## ⇄ ##→meaning) |
| 4 | SHIP | `drain` | broadcast, chain-down copies; re-enters WRITE only via the stamped attest |
| — | PROVE | `qed` | reads every stratum, writes nothing — the CODON byte proves the tower |

> CANONIC INVERTS: the lowercase build commands build the UPPERCASE LANGUAGE of
> CANONICS — the atoms above are that language's words; `nomen` (LANGUAGE's enforcer)
> holds the grammar closed: 8 verbs, byte echo, lexical projection, chain aliases.
