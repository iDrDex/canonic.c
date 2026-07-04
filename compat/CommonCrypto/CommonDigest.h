/* compat/CommonCrypto/CommonDigest.h — the verify-anywhere shim.
 * The kernel is a pure projection and is never edited, so its one borrowed
 * primitive (Apple CommonCrypto SHA256) is mapped to OpenSSL here instead.
 * Used only off-Apple: `make` adds -Icompat and -lcrypto on non-Darwin. */
#pragma once
#include <openssl/sha.h>
typedef SHA256_CTX CC_SHA256_CTX;
typedef unsigned int CC_LONG;
#define CC_SHA256_Init(c)          SHA256_Init(c)
#define CC_SHA256_Update(c, d, n)  SHA256_Update((c), (d), (n))
#define CC_SHA256_Final(md, c)     SHA256_Final((md), (c))
