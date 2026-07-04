/* canonic.h — the public face of the CANONIC KERNEL: three operations, one byte.
 * The implementation (src/canonic.c) is a PURE GOV PROJECTION off the chain and
 * is never edited by hand; this header restates its ABI for linking clients.
 * The byte is the entire interface. Frozen — see PROJECTION. */
#pragma once
#include <stdint.h>
#include <stddef.h>

/* THE 8 BIT-SLOTS — immutable forever; the ABI every client decodes. */
#define GOV_D    0x01u  /* Declarative   — a gov nano carries a named axiom      */
#define GOV_E    0x02u  /* Evidential    — the ledger is fresh (a recent beat)   */
#define GOV_T    0x04u  /* Transparent   — the ledger streams a plan             */
#define GOV_R    0x08u  /* Reproducible  — gov ∧ registry both resolve           */
#define GOV_O    0x10u  /* Operational   — the last event is recent              */
#define GOV_S    0x20u  /* Structural    — chain verifies ∧ category ∈ ontology  */
#define GOV_L    0x40u  /* Linguistic    — chain VALID ∧ learned > 0             */
#define GOV_LANG 0x80u  /* Language      — a nano carries a formal binding       */
#define GOV_255  0xFFu  /* the fixed point — all eight λ-bound                   */

/* OPERATION 1 · SEAL — out = SHA256(prev ++ formal). */
void canonic_seal(const uint8_t prev[32], const char *formal, size_t n,
                  uint8_t out[32]);

/* OPERATION 2 · VERIFY — fold the chain; 1 iff every link reseals. */
typedef struct { const uint8_t (*prev)[32]; const char **formal; const size_t *len;
                 const uint8_t (*hash)[32]; size_t count; } canonic_chain;
int canonic_verify(const canonic_chain *ch);

/* OPERATION 3 · MINT — OR the proven chain-facts into the gov byte. */
typedef struct {
    int has_axiom;      /* D    */
    int fresh;          /* E    */
    int has_plan;       /* T    */
    int gov_and_reg;    /* R    */
    int recent;         /* O    */
    int valid_onto;     /* S    */
    int learned;        /* L    */
    int has_formal;     /* LANG */
} canonic_facts;
uint8_t canonic_gov(const canonic_facts *f);
