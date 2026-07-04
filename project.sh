#!/bin/sh
# project.sh — re-project src/canonic.c from a magic.sov and refresh the purity pin.
# The chain writes the kernel; this script only carries bytes. Usage: ./project.sh <magic.sov>
set -e
SOV="${1:?usage: ./project.sh <magic.sov>}"
tar -xzOf "$SOV" c.chain | python3 -c '
import sys, json
n = json.loads(sys.stdin.readline())
assert n["seq"] == 0, "kernel is seq 0 of c.chain"
sys.stdout.write(n["formal"])' > src/canonic.c
HASH=$(shasum -a 256 src/canonic.c | cut -d" " -f1)
grep '^#' PROJECTION > PROJECTION.tmp
echo "$HASH  src/canonic.c" >> PROJECTION.tmp
mv PROJECTION.tmp PROJECTION
echo "projected: $HASH  src/canonic.c"
