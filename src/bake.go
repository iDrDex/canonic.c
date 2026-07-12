// bake.go — BAKE THE READ PATH (P2). NATIVE CANONIC runtime is Go; .c is bit-mapping only.
// The tsv surfaces (meta/df/term.off/parents) are write-format; a query must never re-parse
// 100MB of text. bake emits three mmap-ready binary sidecars at index time so bm25.c's setup
// is open+mmap+pointer-walk. Each sidecar records the byte length of the tsv prefix it covers
// plus an 8-byte probe of that prefix's last bytes — the reader parses only the tsv DELTA past
// that (the folds since the last bake) and falls back to the full tsv path if the probe
// mismatches (a full rebuild rewrote the surface; the next bake heals it).
//
//	bake <indexdir> <parents.tsv>
//	bake -tomb <indexdir> <parentidx>...   (P3: set tomb=1 on superseded parent records in place)
//
// emits: idx/meta.bin  idx/term.bin  idx/parents.bin   (atomic tmp+mv, published together)
// Layouts are packed little-endian, read by bm25.c — keep the two in sync:
//
//	MetaHdr {u32 'CNM1', u32 nchunks, u64 src, byte probe[8]}          MetaRec {u32 par, u32 dl}
//	TermHdr {u32 'CNT1', u32 nbuckets, u64 src_psl, u64 pool_off}      TermRec {u64 fnv, u64 stroff1, u64 postoff, u64 postlen, u64 df}
//	ParHdr  {u32 'CNP1', u32 nparents, u64 src, byte probe[8], u64 pool_off}
//	ParRec  {byte sha[64], u64 foff, u32 flen, u8 tomb, byte pad[3]}
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

const (
	mMagic = 0x314D4E43 // CNM1
	tMagic = 0x31544E43 // CNT1
	pMagic = 0x31504E43 // CNP1
)

func fnv1a(s []byte) uint64 {
	h := uint64(1469598103934665603)
	for _, c := range s {
		h ^= uint64(c)
		h *= 1099511628211
	}
	return h
}

func probeOf(b []byte) (p [8]byte) {
	n := len(b)
	if n > 8 {
		n = 8
	}
	copy(p[8-n:], b[len(b)-n:])
	return
}

func die(f string, a ...any) { fmt.Fprintf(os.Stderr, "bake: "+f+"\n", a...); os.Exit(69) }

// tsv walks lines yielding raw fields split on the first two tabs (field3 = rest of line)
func lines(b []byte, fn func(f0, f1, f2 []byte)) {
	for len(b) > 0 {
		nl := bytes.IndexByte(b, '\n')
		var ln []byte
		if nl < 0 {
			ln, b = b, nil
		} else {
			ln, b = b[:nl], b[nl+1:]
		}
		if len(ln) == 0 {
			continue
		}
		t1 := bytes.IndexByte(ln, '\t')
		if t1 < 0 {
			fn(ln, nil, nil)
			continue
		}
		rest := ln[t1+1:]
		t2 := bytes.IndexByte(rest, '\t')
		if t2 < 0 {
			fn(ln[:t1], rest, nil)
			continue
		}
		fn(ln[:t1], rest[:t2], rest[t2+1:])
	}
}

func atol(b []byte) int64 { v, _ := strconv.ParseInt(string(b), 10, 64); return v }

func writeAtomic(path string, emit func(w *bufio.Writer) error) {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		die("%v", err)
	}
	w := bufio.NewWriterSize(f, 1<<20)
	if err := emit(w); err != nil {
		die("%v", err)
	}
	w.Flush()
	f.Close()
	if err := os.Rename(tmp, path); err != nil {
		die("%v", err)
	}
}

func tomb(dir string, idxs []string) {
	// P3: poke tomb=1 on parent records in place — an 80-byte record per parent, tomb at +76.
	const hdr = 32 // ParHdr
	const rec = 80
	f, err := os.OpenFile(dir+"/parents.bin", os.O_RDWR, 0)
	if err != nil {
		die("%v", err)
	}
	defer f.Close()
	var h [8]byte
	if _, err := f.ReadAt(h[:], 0); err != nil {
		die("%v", err)
	}
	if binary.LittleEndian.Uint32(h[:4]) != pMagic {
		die("parents.bin: bad magic")
	}
	np := int64(binary.LittleEndian.Uint32(h[4:8]))
	n := 0
	for _, s := range idxs {
		p, _ := strconv.ParseInt(s, 10, 64)
		if p < 1 || p > np {
			continue
		}
		if _, err := f.WriteAt([]byte{1}, hdr+(p-1)*rec+76); err != nil {
			die("%v", err)
		}
		n++
	}
	fmt.Printf("bake · tombed %d parent(s)\n", n)
}

func main() {
	args := os.Args[1:]
	if len(args) >= 2 && args[0] == "-tomb" {
		tomb(args[1], args[2:])
		return
	}
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "bake <indexdir> <parents.tsv> | bake -tomb <indexdir> <parentidx>...")
		os.Exit(64)
	}
	dir, pf := args[0], args[1]

	// ---- meta.bin : chunkid -> {parent, dl} (record i = chunk i, 1-based like meta.tsv) ----
	mt, err := os.ReadFile(dir + "/meta.tsv")
	if err != nil {
		die("no meta.tsv")
	}
	maxc := int64(0)
	lines(mt, func(f0, f1, f2 []byte) {
		if cl := atol(f0); cl > maxc {
			maxc = cl
		}
	})
	mr := make([]uint64, maxc+1) // par<<32|dl packed for the walk, unpacked at write
	lines(mt, func(f0, f1, f2 []byte) {
		cl := atol(f0)
		if cl > 0 && cl <= maxc {
			mr[cl] = uint64(atol(f1))<<32 | uint64(uint32(atol(f2)))
		}
	})
	mprobe := probeOf(mt)
	writeAtomic(dir+"/meta.bin", func(w *bufio.Writer) error {
		binary.Write(w, binary.LittleEndian, uint32(mMagic))
		binary.Write(w, binary.LittleEndian, uint32(maxc))
		binary.Write(w, binary.LittleEndian, uint64(len(mt)))
		w.Write(mprobe[:])
		var b [8]byte
		for c := int64(1); c <= maxc; c++ {
			binary.LittleEndian.PutUint32(b[:4], uint32(mr[c]>>32))
			binary.LittleEndian.PutUint32(b[4:], uint32(mr[c]))
			w.Write(b[:])
		}
		return nil
	})
	mr = nil

	// ---- term.bin : fnv open-addressing -> {stroff1, postoff, postlen, df} + string pool ----
	to, err := os.ReadFile(dir + "/term.off")
	if err != nil {
		die("no term.off")
	}
	dfm := map[uint64]uint64{}
	if dfb, err := os.ReadFile(dir + "/df.tsv"); err == nil {
		lines(dfb, func(f0, f1, f2 []byte) { dfm[fnv1a(f0)] = uint64(atol(f1)) })
	}
	nterm := int64(0)
	lines(to, func(f0, f1, f2 []byte) { nterm++ })
	nb := uint32(1)
	for int64(nb) < nterm*2+16 {
		nb <<= 1
	}
	type trec struct{ fnv, stroff1, postoff, postlen, df uint64 }
	tr := make([]trec, nb)
	pool := make([]byte, 0, len(to))
	lines(to, func(f0, f1, f2 []byte) {
		h := fnv1a(f0)
		i := uint32(h) & (nb - 1)
		for tr[i].stroff1 != 0 {
			if tr[i].fnv == h {
				return
			}
			i = (i + 1) & (nb - 1)
		}
		df := uint64(1)
		if v, ok := dfm[h]; ok {
			df = v
		}
		tr[i] = trec{h, uint64(len(pool)) + 1, uint64(atol(f1)), uint64(atol(f2)), df}
		pool = append(pool, f0...)
		pool = append(pool, 0)
	})
	var psl uint64
	if st, err := os.Stat(dir + "/post.sorted"); err == nil {
		psl = uint64(st.Size())
	}
	writeAtomic(dir+"/term.bin", func(w *bufio.Writer) error {
		binary.Write(w, binary.LittleEndian, uint32(tMagic))
		binary.Write(w, binary.LittleEndian, nb)
		binary.Write(w, binary.LittleEndian, psl)
		binary.Write(w, binary.LittleEndian, uint64(24)+uint64(nb)*40)
		var b [40]byte
		for _, r := range tr {
			binary.LittleEndian.PutUint64(b[0:], r.fnv)
			binary.LittleEndian.PutUint64(b[8:], r.stroff1)
			binary.LittleEndian.PutUint64(b[16:], r.postoff)
			binary.LittleEndian.PutUint64(b[24:], r.postlen)
			binary.LittleEndian.PutUint64(b[32:], r.df)
			w.Write(b[:])
		}
		w.Write(pool)
		return nil
	})
	tr = nil

	// ---- parents.bin : line idx (1-based) -> {sha[64], formal off/len, tomb} + formal pool ----
	pb, err := os.ReadFile(pf)
	if err != nil {
		die("no parents")
	}
	type prec struct {
		sha  []byte
		foff uint64
		flen uint32
	}
	var prs []prec
	var fpn uint64
	// formal = everything after the FIRST tab (a formal may itself carry escaped tabs);
	// pool offsets accumulate here, the pool bytes stream in write pass 2 — zero copies held.
	{
		b := pb
		for len(b) > 0 {
			nl := bytes.IndexByte(b, '\n')
			var ln []byte
			if nl < 0 {
				ln, b = b, nil
			} else {
				ln, b = b[:nl], b[nl+1:]
			}
			if len(ln) == 0 {
				continue
			}
			t1 := bytes.IndexByte(ln, '\t')
			r := prec{}
			var hash, formal []byte
			if t1 < 0 {
				hash = ln
			} else {
				hash, formal = ln[:t1], ln[t1+1:]
			}
			if len(hash) == 64 {
				r.sha = hash
			}
			if len(formal) > 0 {
				r.foff = fpn
				r.flen = uint32(len(formal))
				fpn += uint64(len(formal))
			}
			prs = append(prs, r)
		}
	}
	// P3: dir/tombs = 1-based parent line indices of SUPERSEDED named folds (fossils). The chain
	// keeps every fold (append-only); the DERIVED read path answers with the latest register only.
	// tombs is regenerated by a full build and appended by the inc — bake re-applies it every run.
	tombed := map[int64]bool{}
	if tb, err := os.ReadFile(dir + "/tombs"); err == nil {
		lines(tb, func(f0, f1, f2 []byte) {
			if v := atol(f0); v > 0 {
				tombed[v] = true
			}
		})
	}
	pprobe := probeOf(pb)
	writeAtomic(dir+"/parents.bin", func(w *bufio.Writer) error {
		binary.Write(w, binary.LittleEndian, uint32(pMagic))
		binary.Write(w, binary.LittleEndian, uint32(len(prs)))
		binary.Write(w, binary.LittleEndian, uint64(len(pb)))
		w.Write(pprobe[:])
		binary.Write(w, binary.LittleEndian, uint64(32)+uint64(len(prs))*80)
		var rec [80]byte
		for pi, r := range prs {
			for i := range rec {
				rec[i] = 0
			}
			copy(rec[:64], r.sha)
			binary.LittleEndian.PutUint64(rec[64:], r.foff)
			binary.LittleEndian.PutUint32(rec[72:], r.flen)
			if tombed[int64(pi)+1] {
				rec[76] = 1
			}
			w.Write(rec[:])
		}
		// pool pass: formals in line order
		b := pb
		for len(b) > 0 {
			nl := bytes.IndexByte(b, '\n')
			var ln []byte
			if nl < 0 {
				ln, b = b, nil
			} else {
				ln, b = b[:nl], b[nl+1:]
			}
			if len(ln) == 0 {
				continue
			}
			if t1 := bytes.IndexByte(ln, '\t'); t1 >= 0 && t1+1 < len(ln) {
				w.Write(ln[t1+1:])
			}
		}
		return nil
	})
	fmt.Printf("bake · chunks=%d terms=%d parents=%d pools=%.1fMB\n", maxc, nterm, len(prs), float64(uint64(len(pool))+fpn)/1048576.0)
}
