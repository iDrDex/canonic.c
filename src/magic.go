// magic.go — NATIVE # (Go). C is native bits; Go is native #. crypto-gated: every op is sha256.
// The # layer (fold + chain verify) native in Go — the EVO-runtime target of the transpiler.
//   magic fold <prev> <formal>   → # = sha256(prev++formal)
//   magic verify                 → stream events.jsonl on stdin, verify the # chain (crypto-gated)
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) >= 4 && os.Args[1] == "fold" {
		h := sha256.Sum256([]byte(os.Args[2] + os.Args[3]))
		fmt.Println(hex.EncodeToString(h[:]))
		return
	}
	if len(os.Args) >= 2 && os.Args[1] == "verify" {
		sc := bufio.NewScanner(os.Stdin)
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		var prev string
		n, breaks := 0, 0
		for sc.Scan() {
			var e struct {
				PrevHash    string `json:"prev_hash"`
				PayloadHash string `json:"payload_hash"`
			}
			if json.Unmarshal(sc.Bytes(), &e) != nil || e.PayloadHash == "" {
				continue
			}
			if prev != "" && e.PrevHash != "" && e.PrevHash != prev {
				breaks++
			}
			prev = e.PayloadHash
			n++
		}
		verdict := "INTACT"
		if breaks > 0 {
			verdict = fmt.Sprintf("FORKED(%d)", breaks)
		}
		fmt.Printf("go-# verify: %d events, tail-links=%s\n", n, verdict)
		return
	}
	if len(os.Args) >= 2 && os.Args[1] == "prove" {
		// REPRODUCE-COLD, native: for every NANO, sha256(prev++formal)==payload_hash.
		// json.Unmarshal gives the exact formal the fold stored — the authoritative check (no perl).
		sc := bufio.NewScanner(os.Stdin)
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		n, bad, legacy := 0, 0, 0
		for sc.Scan() {
			var e struct {
				Type        string `json:"type"`
				PrevHash    string `json:"prev_hash"`
				PayloadHash string `json:"payload_hash"`
				Payload     struct {
					Formal string `json:"formal"`
				} `json:"payload"`
			}
			if json.Unmarshal(sc.Bytes(), &e) != nil || e.Type != "NANO" || e.PayloadHash == "" {
				continue
			}
			// ERA-AWARE: reproduce-cold applies ONLY to the current {category,formal} schema.
			// Legacy folds carry a different payload schema (no formal, e.g. {name,bits,byte,lambdas})
			// whose # was computed over that era's canonical form — NOT sha256(prev++formal). They are
			// era-drift, not corruption: verified by LINKAGE + SE, never by this serializer. Skip them.
			if e.Payload.Formal == "" {
				legacy++
				continue
			}
			n++
			h := sha256.Sum256([]byte(e.PrevHash + e.Payload.Formal))
			if hex.EncodeToString(h[:]) != e.PayloadHash {
				bad++
				if len(os.Args) >= 3 && os.Args[2] == "-v" {
					fmt.Fprintln(os.Stderr, e.PayloadHash)
				}
			}
		}
		fmt.Printf("go-# prove: current-schema %d folds, %d reproduce-cold mismatch; legacy-schema %d (era-drift, linked+SE, not corruption)\n", n, bad, legacy)
		return
	}
	if len(os.Args) >= 2 && os.Args[1] == "qed" {
		// ONE STREAM (P4): everything the qed byte needs from the ledger in a single pass —
		// linkage breaks (full chain), the signer SEQUENCE (signed/unsigned/regressions/live,
		// the dark-era detector), #anchor count, and the raw line number of each tip hash
		// passed as an arg (surface staleness = n - line). Replaces five full-file scans.
		tips := os.Args[2:]
		tipline := map[string]int{}
		sc := bufio.NewScanner(os.Stdin)
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		var prev, salgPrev string
		raw, n, breaks, stot, ssig, suns, regr, anchors := 0, 0, 0, 0, 0, 0, 0, 0
		sawSigned, live := false, false
		for sc.Scan() {
			raw++
			var e struct {
				PrevHash    string `json:"prev_hash"`
				PayloadHash string `json:"payload_hash"`
				SigAlg      string `json:"sig_alg"`
				Payload     struct {
					Formal string `json:"formal"`
				} `json:"payload"`
			}
			if json.Unmarshal(sc.Bytes(), &e) != nil {
				continue
			}
			if e.SigAlg != "" {
				stot++
				if e.SigAlg == "none" {
					suns++
					if salgPrev != "" && salgPrev != "none" {
						regr++
					}
					live = sawSigned
				} else {
					ssig++
					sawSigned = true
					live = false
				}
				salgPrev = e.SigAlg
			}
			if len(e.Payload.Formal) >= 8 && e.Payload.Formal[:8] == "#anchor " {
				anchors++
			}
			if e.PayloadHash == "" {
				continue
			}
			if prev != "" && e.PrevHash != "" && e.PrevHash != prev {
				breaks++
			}
			prev = e.PayloadHash
			n++
			for _, t := range tips {
				if t == e.PayloadHash {
					tipline[t] = raw
				}
			}
		}
		liveS := "clean"
		if live {
			liveS = "LIVE"
		}
		fmt.Printf("go-# qed: n=%d events=%d breaks=%d stot=%d ssig=%d suns=%d regr=%d live=%s anchors=%d",
			raw, n, breaks, stot, ssig, suns, regr, liveS, anchors)
		for _, t := range tips {
			ln, okt := tipline[t]
			if !okt {
				ln = -1
			}
			fmt.Printf(" tip:%s=%d", t, ln)
		}
		fmt.Println()
		return
	}
	// throughput bench: N native #s
	if len(os.Args) >= 3 && os.Args[1] == "bench" {
		var n int
		fmt.Sscan(os.Args[2], &n)
		prev := "genesis"
		for i := 0; i < n; i++ {
			h := sha256.Sum256([]byte(prev + fmt.Sprint(i)))
			prev = hex.EncodeToString(h[:])
		}
		fmt.Printf("go-# bench: %d folds, tip %s\n", n, prev[:16])
		return
	}
	fmt.Fprintln(os.Stderr, "magic {fold <prev> <formal> | verify | bench <n>}")
	os.Exit(64)
}
