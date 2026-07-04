/* ═══════════════════════════════════════════════════════════════════════════
 * canonic.c — THE CANONIC KERNEL.  It BOOTS Magic OS.  The truly minimal bitwise
 * GOV encoder: three operations, one byte.  CANONIC is the kernel (the paradigm,
 * the outside looking in, ∩); MAGIC is the OS it boots (magic.c the runtime, magic.go
 * the evolution runtime).  This is the easter egg: canonic.c rides inside magic.sov
 * as one nano — open GOV without a source repo.  Decode, concat, and you hold the
 * kernel that boots the OS.  canonic.c boots magic; the split is CANONIC / MAGIC.
 *
 * THE BOUNDARY: the CANONIC KERNEL is PURE BITWISE ENCODING OF GOV.  Three
 * operations and one byte.  It never stats the tree; it reads only the chain's
 * bits.  MAGIC — the runtime it boots (walk · migrate · embed · the surface) —
 * lives above, in magic.c and magic.go.  CANONIC boots; MAGIC runs.
 * The test is one line: touches a file → runtime; touches only bits → kernel.
 *
 * THE 8-BIT ABI — the REAL compiled boundary (these comments compile away; the
 * bit-assignments do NOT).  The byte crossing this line IS the entire public
 * interface.  Every client on every platform decodes against exactly this:
 *
 *   D    0x01  Declarative   a gov nano carries a named axiom
 *   E    0x02  Evidential    the ledger is fresh (a recent beat)
 *   T    0x04  Transparent   the ledger streams a plan (the road is visible)
 *   R    0x08  Reproducible  gov ∧ registry both resolve (the seed regrows)
 *   O    0x10  Operational   the last event is recent (the walk is live)
 *   S    0x20  Structural    the chain verifies ∧ a category is in the ontology
 *   L    0x40  Linguistic    the chain is VALID ∧ learned > 0 (it has learned)
 *   LANG 0x80  Language      a nano carries a "formal" natural-language binding
 *
 * min-LOG ≡ max-GOV ≡ 255 = 0xFF.  The fixed point is a full byte; you cannot
 * remove a bit from a contract that IS the fixed point.  Frozen.
 * ═════════════════════════════════════════════════════════════════════════ */

#include <stdint.h>
#include <stddef.h>
#include <string.h>
#include <CommonCrypto/CommonDigest.h>   /* the ONE borrowed primitive: SHA256 */

/* ── THE 8 BIT-SLOTS — immutable forever; the ABI every client decodes ────── */
#define GOV_D 0x01u
#define GOV_E 0x02u
#define GOV_T 0x04u
#define GOV_R 0x08u
#define GOV_O 0x10u
#define GOV_S 0x20u
#define GOV_L 0x40u
#define GOV_LANG 0x80u
#define GOV_255 0xFFu   /* the fixed point — all eight λ-bound */

/* ── OPERATION 1 · SEAL — hash(prev ++ content) → 32 bytes.  A nano's identity
 *    is the SHA256 of its predecessor's hash concatenated with its own formal.
 *    This is the whole protocol: the raw hash of the next. ─────────────────── */
void canonic_seal(const uint8_t prev[32], const char *formal, size_t n,
                 uint8_t out[32]) {
    CC_SHA256_CTX c;
    CC_SHA256_Init(&c);
    CC_SHA256_Update(&c, prev, 32);
    CC_SHA256_Update(&c, formal, (CC_LONG)n);
    CC_SHA256_Final(out, &c);
}

/* ── OPERATION 2 · VERIFY — fold the chain link by link.  Recompute each hash
 *    from the running prev; any mismatch breaks the seal of all that follows.
 *    Returns 1 (VALID) iff the recomputed head equals the stored head. ─────── */
typedef struct { const uint8_t (*prev)[32]; const char **formal; const size_t *len;
                 const uint8_t (*hash)[32]; size_t count; } canonic_chain;

int canonic_verify(const canonic_chain *ch) {
    uint8_t run[32] = {0};                    /* genesis: prev = 32 zero bytes */
    for (size_t i = 0; i < ch->count; i++) {
        uint8_t h[32];
        if (memcmp(ch->prev[i], run, 32) != 0) return 0;      /* prev-link broke */
        canonic_seal(run, ch->formal[i], ch->len[i], h);
        if (memcmp(h, ch->hash[i], 32) != 0) return 0;        /* content broke   */
        memcpy(run, h, 32);                                   /* advance the head*/
    }
    return 1;
}

/* ── OPERATION 3 · MINT — the 8-bit gov byte from CHAIN BITS ONLY.  Each λ is a
 *    pure predicate over facts already proven above (the verify result, counts,
 *    flags carried in the chain).  No file is stat'd; no tree is walked.  The
 *    caller passes the proven bits; the kernel ORs them into the contract. ─── */
typedef struct {
    int has_axiom;      /* D: a gov nano names an axiom          */
    int fresh;          /* E: a recent beat is on the chain      */
    int has_plan;       /* T: a WALK_PLAN streams                 */
    int gov_and_reg;    /* R: gov ∧ registry both resolve        */
    int recent;         /* O: last event is recent               */
    int valid_onto;     /* S: chain valid ∧ category ∈ ontology  */
    int learned;        /* L: chain VALID ∧ learned > 0          */
    int has_formal;     /* LANG: a nano carries a formal binding */
} canonic_facts;

uint8_t canonic_gov(const canonic_facts *f) {
    uint8_t b = 0;
    if (f->has_axiom)   b |= GOV_D;
    if (f->fresh)       b |= GOV_E;
    if (f->has_plan)    b |= GOV_T;
    if (f->gov_and_reg) b |= GOV_R;
    if (f->recent)      b |= GOV_O;
    if (f->valid_onto)  b |= GOV_S;
    if (f->learned)     b |= GOV_L;
    if (f->has_formal)  b |= GOV_LANG;
    return b;                                 /* the byte IS the boundary (255) */
}

/* THE WHOLE KERNEL: seal ∘ verify ∘ gov.  Three functions, one byte, one hash.
 * Everything above this line is the runtime; everything here is truth in bits.
 * The comments are the easter egg — they name the bits, then compile away, and
 * what remains on the wire is the pure bitwise encoding of GOV.  QED. */
