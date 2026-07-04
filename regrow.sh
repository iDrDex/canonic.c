#!/bin/sh
# regrow.sh — THE PROJECTION GATE. canonic.sov is the artifact (sov code); file
# trees are its projections. Two modes:
#   ./regrow.sh              gate: this tree must be EXACTLY the sov's regrowth
#   ./regrow.sh --into DIR   materialize: regrow the projection from the sov alone
# The chain verifies by the kernel's own law: SHA256(prev32 ++ formal), genesis
# 32 zero bytes — canonic_verify semantics. Change flows sov → tree, never back.
set -e
DIR=""
[ "$1" = "--into" ] && DIR="${2:?usage: ./regrow.sh --into DIR}"
python3 - "$DIR" <<'EOF'
import hashlib, json, tarfile, sys, os
into = sys.argv[1]
t = tarfile.open('canonic.sov')
man = json.loads(t.extractfile('sov.manifest').read())
nanos = [json.loads(l) for l in t.extractfile('canonic.chain').read().decode().splitlines()]
run = b'\x00' * 32
for n in nanos:
    assert bytes.fromhex(n['prev']) == run, f"chain broke at seq {n['seq']}"
    h = hashlib.sha256(run + n['formal'].encode()).digest()
    assert h.hex() == n['hash'], f"seal broke at seq {n['seq']}"
    run = h
assert run.hex() == man['head'], 'head mismatch'
print(f"sov chain: VALID · {len(nanos)} nanos · head {man['head'][:16]}")
if into:
    for n in nanos:
        p = os.path.join(into, n['path'])
        os.makedirs(os.path.dirname(p) or '.', exist_ok=True)
        open(p, 'w', encoding='utf-8').write(n['formal'])
        if p.endswith('.sh'):
            os.chmod(p, 0o755)
    print(f"projection: REGROWN into {into}/ ({len(nanos)} files)")
else:
    drift = [n['path'] for n in nanos
             if open(n['path'], encoding='utf-8').read() != n['formal']]
    assert not drift, f"TREE DRIFTED FROM THE SOV: {drift}"
    print("projection: EXACT — the tree is the sov's regrowth")
EOF
