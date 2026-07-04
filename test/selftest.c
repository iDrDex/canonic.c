/* test/selftest.c — the e2e proof, exactly the kernel's own claim:
 * chain VALID · tampered BROKEN (caught) · gov 255/255 · decayed 237/255.
 * Exit 0 iff every gate holds. */
#include <stdio.h>
#include <string.h>
#include "canonic.h"

int main(void) {
    const char *formal[3] = {
        "AXIOM: INTEL = COIN = WORK = LEARNING (A1)",
        "the chain is the term sheet",
        "min-LOG == max-GOV == 255"
    };
    size_t len[3];
    uint8_t prev[3][32], hash[3][32], run[32] = {0};
    for (int i = 0; i < 3; i++) {
        len[i] = strlen(formal[i]);
        memcpy(prev[i], run, 32);
        canonic_seal(run, formal[i], len[i], hash[i]);
        memcpy(run, hash[i], 32);
    }

    canonic_chain ch = { (const uint8_t (*)[32])prev, formal, len,
                         (const uint8_t (*)[32])hash, 3 };
    int ok = canonic_verify(&ch);
    printf("chain: %s\n", ok ? "VALID" : "BROKEN");
    int pass = ok;

    const char *tampered = "the chain is the term sheet (edited)";
    const char *f2[3] = { formal[0], tampered, formal[2] };
    size_t l2[3] = { len[0], strlen(tampered), len[2] };
    canonic_chain bad = { (const uint8_t (*)[32])prev, f2, l2,
                          (const uint8_t (*)[32])hash, 3 };
    int caught = !canonic_verify(&bad);
    printf("tampered: %s\n", caught ? "BROKEN (caught)" : "MISSED");
    pass = pass && caught;

    canonic_facts all = {1, 1, 1, 1, 1, 1, 1, 1};
    uint8_t b = canonic_gov(&all);
    printf("gov: %u/255\n", b);
    pass = pass && (b == GOV_255);

    canonic_facts decayed = all;
    decayed.fresh = 0;
    decayed.recent = 0;
    uint8_t d = canonic_gov(&decayed);
    printf("decayed: %u/255 (E,O unbound)\n", d);
    pass = pass && (d == (uint8_t)(GOV_255 & (uint8_t)~(GOV_E | GOV_O)));

    puts(pass ? "ALL GATES PASS" : "GATE FAILED");
    return pass ? 0 : 1;
}
