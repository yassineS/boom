// Copyright ©2012 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package boom is a wrapper for the htslib library.
package boom

// https://www.htslib.org/

/*
#cgo CFLAGS: -g -O2 -fPIC -m64 -pthread
#cgo LDFLAGS: -lhts
#include <htslib/sam.h>
#include <htslib/hts.h>
void setBin(bam1_t *b, uint16_t bin)        { b->core.bin = bin; }
void setQual(bam1_t *b, uint8_t flag)       { b->core.flag = flag; }
void setLQname(bam1_t *b, uint8_t l_qname)  { b->core.l_qname = l_qname; }
void setFlag(bam1_t *b, uint16_t flag)      { b->core.flag = flag; }
void setNCigar(bam1_t *b, uint16_t n_cigar) { b->core.n_cigar = n_cigar; }
*/
import "C"

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"unsafe"
)

var (
	valueIsNil       = fmt.Errorf("boom: value is nil")
	notBamFile       = fmt.Errorf("boom: not bam file")
	couldNotAllocate = fmt.Errorf("boom: could not allocate")
	cannotAddr       = fmt.Errorf("boom: cannot address value")
	bamIsBigEndian   = C.ed_is_big() == 1
	endian           = [2]binary.ByteOrder{
		binary.LittleEndian,
		binary.BigEndian,
	}[C.ed_is_big()]
)

var (
	noHeader = errors.New("boom: no header")
)

// Verbosity sets and returns the level of debugging information emitted on stderr by htslib.
// The level of verbosity interpreted by htslib ranges from 0 (HTS_LOG_OFF) to 5 (HTS_LOG_TRACE) inclusive.
// Lower values being less verbose. Passing values of v outside the range [0, 5] do not alter verbosity.
func Verbosity(v int) int {
	if 0 <= v && v <= 5 {
		C.hts_set_log_level((C.enum_htsLogLevel)(v))
	}
	return int(C.hts_get_log_level())
}

// A bamRecord wraps the bam1_t BAM record.
type bamRecord struct {
	b *C.bam1_t
}

// newBamRecord creates a new bamRecord wrapping b or a newly malloc'd bam1_t if b is nil,
// and setting a finaliser that C.free()s the contained bam1_t.
// newBamRecord should always be used unless the bamRecord will be explicitly memory managed, or
// wraps a bam1_t that will be memory managed elsewhere.
func newBamRecord(b *C.bam1_t) (br *bamRecord, err error) {
	if b == nil {
		b = (*C.bam1_t)(unsafe.Pointer(C.malloc((C.size_t)(unsafe.Sizeof(C.bam1_t{})))))

		if b == nil {
			return nil, couldNotAllocate
		}
		*b = C.bam1_t{}
	}

	br = &bamRecord{b}
	runtime.SetFinalizer(br, (*bamRecord).bamRecordFree)

	return
}

// The following methods are helpers to safely return bam1_t field values.
// All first check that the pointer to the bam1_t is not nil and convert to the appropriate
// Go type.
func (br *bamRecord) tid() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int32(br.b.core.tid)
}
func (br *bamRecord) setTid(tid int32) {
	if br.b == nil {
		panic(valueIsNil)
	}
	br.b.core.tid = C.int32_t(tid)
}
func (br *bamRecord) pos() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int32(br.b.core.pos)
}
func (br *bamRecord) setPos(pos int32) {
	if br.b == nil {
		panic(valueIsNil)
	}
	br.b.core.pos = C.hts_pos_t(pos)
}
func (br *bamRecord) bin() uint16 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return uint16(br.b.core.bin)
}
func (br *bamRecord) setBin(bin uint16) {
	if br.b == nil {
		panic(valueIsNil)
	}
	C.setBin(br.b, C.uint16_t(bin))
}
func (br *bamRecord) qual() byte {
	if br.b == nil {
		panic(valueIsNil)
	}
	return byte(br.b.core.qual)
}
func (br *bamRecord) setQual(qual byte) {
	if br.b == nil {
		panic(valueIsNil)
	}
	C.setQual(br.b, C.uint8_t(qual))
}
func (br *bamRecord) lQname() byte {
	if br.b == nil {
		panic(valueIsNil)
	}
	return byte(br.b.core.l_qname)
}
func (br *bamRecord) setLQname(lQname byte) {
	if br.b == nil {
		panic(valueIsNil)
	}
	C.setLQname(br.b, C.uint8_t(lQname))
}
func (br *bamRecord) flag() Flags {
	if br.b == nil {
		panic(valueIsNil)
	}
	return Flags(br.b.core.flag)
}
func (br *bamRecord) setFlag(flags Flags) {
	if br.b == nil {
		panic(valueIsNil)
	}
	C.setFlag(br.b, C.uint16_t(flags))
}
func (br *bamRecord) nCigar() uint16 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return uint16(br.b.core.n_cigar)
}
func (br *bamRecord) setNCigar(nCigar uint16) {
	if br.b == nil {
		panic(valueIsNil)
	}
	C.setNCigar(br.b, C.uint16_t(nCigar))
}
func (br *bamRecord) lQseq() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int32(br.b.core.l_qseq)
}
func (br *bamRecord) setLQseq(lQseq int32) {
	if br.b == nil {
		panic(valueIsNil)
	}
	br.b.core.l_qseq = C.int32_t(lQseq)
}
func (br *bamRecord) mtid() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int32(br.b.core.mtid)
}
func (br *bamRecord) setMtid() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int32(br.b.core.mtid)
}
func (br *bamRecord) mpos() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int32(br.b.core.mpos)
}
func (br *bamRecord) setMpos(mpos int32) {
	if br.b == nil {
		panic(valueIsNil)
	}
	br.b.core.mpos = C.hts_pos_t(mpos)
}
func (br *bamRecord) isize() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int32(br.b.core.isize)
}
func (br *bamRecord) setIsize(isize int32) {
	if br.b == nil {
		panic(valueIsNil)
	}
	br.b.core.isize = C.hts_pos_t(isize)
}
func (br *bamRecord) lAux() int32 {
	if br.b == nil {
		panic(valueIsNil)
	}
	// Calculate aux length using the formula from bam_get_l_aux macro
	// l_data - (n_cigar<<2) - l_qname - l_qseq - ((l_qseq + 1)>>1)
	return int32(br.b.l_data) - int32(br.b.core.n_cigar<<2) - int32(br.b.core.l_qname) - int32(br.b.core.l_qseq) - int32((br.b.core.l_qseq+1)>>1)
}
func (br *bamRecord) setLAux(lAux int32) {
	// This is now a calculated field, so we can't set it directly
	// This function is kept for compatibility but does nothing
	if br.b == nil {
		panic(valueIsNil)
	}
}
func (br *bamRecord) dataLen() int {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int(br.b.l_data)
}
func (br *bamRecord) dataCap() int {
	if br.b == nil {
		panic(valueIsNil)
	}
	return int(br.b.m_data)
}
func (br *bamRecord) dataPtr() uintptr {
	if br.b == nil {
		panic(valueIsNil)
	}
	return uintptr(unsafe.Pointer(br.b.data))
}
func (br *bamRecord) dataUnsafe() []byte {
	if br.b == nil {
		panic(valueIsNil)
	}

	l := int(br.b.l_data)
	var data []byte
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sliceHeader.Cap = l
	sliceHeader.Len = l
	sliceHeader.Data = uintptr(unsafe.Pointer(br.b.data))

	return data
}
func (br *bamRecord) setDataUnsafe(data []byte) {
	if br.b == nil {
		panic(valueIsNil)
	}

	l := len(data)
	if br.dataCap() < len(data) {
		if br.b.data == nil {
			br.b.data = (*C.uint8_t)(unsafe.Pointer(C.malloc((C.size_t)(l))))
		} else {

			br.b.data = (*C.uint8_t)(unsafe.Pointer(C.realloc(unsafe.Pointer(br.b.data), (C.size_t)(l))))
		}
		if br.b.data == nil {
			panic(couldNotAllocate)
		}
	}

	var newData []byte
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&newData))
	sliceHeader.Cap = l
	sliceHeader.Len = l
	sliceHeader.Data = uintptr(unsafe.Pointer(br.b.data))
	copy(newData, data)
}

// bamRecordFree C.free()s the contained bam1_t and its data, first checking for nil pointers.
func (br *bamRecord) bamRecordFree() {
	if br.b != nil {
		if br.b.data != nil {
			C.free(unsafe.Pointer(br.b.data))
		}
		C.free(unsafe.Pointer(br.b))
		br.b = nil
	}
}

// A samFile wraps a samFile (htsFile).
type samFile struct {
	fp  *C.samFile
	hdr *C.sam_hdr_t
}

// samOpen/samFdOpen open a SAM or BAM file with the given filename/fd, mode and optional auxilliary header.
// According to htslib:
//
// mode matches /[rwa](b?)(c?)(g?)(u?)(z?)[0-9]?/
//
//   'r' for reading,
//   'w' for writing,
//   'a' for appending,
//   'b' for BAM format,
//   'c' for CRAM format,
//   'g' for gzip compression,
//   'u' for uncompressed,
//   'z' for bgzf compression,
//   '0'-'9' for compression level
//
// If mode[0] == 'w' or 'a', aux should be a bamHeader containing header information to write.
// For reading, the header will be read from the file.
func samOpen(filename, mode string, aux header) (sf *samFile, err error) {
	fn, m := C.CString(filename), C.CString(mode)
	defer C.free(unsafe.Pointer(fn))
	defer C.free(unsafe.Pointer(m))

	fp := C.hts_open(
		(*C.char)(unsafe.Pointer(fn)),
		(*C.char)(unsafe.Pointer(m)),
	)
	if fp == nil {
		return nil, fmt.Errorf("boom: failed to open file %s", filename)
	}

	var hdr *C.sam_hdr_t
	if mode[0] == 'r' {
		// Read the header from the file
		hdr = C.sam_hdr_read(fp)
		if hdr == nil {
			C.hts_close(fp)
			return nil, fmt.Errorf("boom: failed to read header from %s", filename)
		}
	} else {
		// Writing mode - use provided header
		switch a := aux.(type) {
		case *bamHeader:
			if a != nil && a.bh != nil {
				hdr = C.sam_hdr_dup(a.bh)
				if hdr == nil {
					C.hts_close(fp)
					return nil, fmt.Errorf("boom: failed to duplicate header")
				}
				// Write the header to the file
				if C.sam_hdr_write(fp, hdr) < 0 {
					C.sam_hdr_destroy(hdr)
					C.hts_close(fp)
					return nil, fmt.Errorf("boom: failed to write header")
				}
			} else {
				C.hts_close(fp)
				return nil, noHeader
			}
		case nil:
			C.hts_close(fp)
			return nil, noHeader
		default:
			C.hts_close(fp)
			return nil, fmt.Errorf("boom: unsupported header type %T", aux)
		}
	}

	sf = &samFile{fp: fp, hdr: hdr}
	runtime.SetFinalizer(sf, (*samFile).samClose)

	return
}
func samFdOpen(fd uintptr, mode string, aux header) (sf *samFile, err error) {
	// Note: htslib doesn't have a direct fdopen equivalent for samFile
	// We'll need to use hts_hopen with an hFILE wrapper, but that's complex
	// For now, we'll return an error indicating this is not supported
	return nil, fmt.Errorf("boom: samFdOpen not supported with htslib - use samOpen instead")
}

// header returns the bamHeader wrapping the sam_hdr_t associated with sf
func (sf *samFile) header() *bamHeader {
	if sf.hdr == nil {
		return nil
	}
	return &bamHeader{bh: sf.hdr}
}

// samClose closes the samFile, freeing the C data allocations.
func (sf *samFile) samClose() error {
	if sf.fp == nil {
		return valueIsNil
	}
	runtime.SetFinalizer(sf, nil)

	if sf.hdr != nil {
		C.sam_hdr_destroy(sf.hdr)
		sf.hdr = nil
	}

	C.hts_close(sf.fp)
	sf.fp = nil

	return nil
}

// samRead reads and returns the next BAM record returning the number of bytes read,
// a *bamRecord containing the record data and any error that occurred.
func (sf *samFile) samRead() (n int, br *bamRecord, err error) {
	if sf.fp == nil || sf.hdr == nil {
		return 0, nil, valueIsNil
	}

	br, err = newBamRecord(nil)
	if err != nil {
		return
	}

	cn := C.sam_read1(
		sf.fp,
		sf.hdr,
		br.b,
	)
	n = int(cn)
	if n < 0 {
		err = io.EOF
	}

	return
}

// samWrite writes a BAM record represented by br, returning the number of bytes written
// and any error that occurred.
func (sf *samFile) samWrite(br *bamRecord) (n int, err error) {
	if sf.fp == nil || sf.hdr == nil || br.b == nil {
		return 0, valueIsNil
	}

	ret := C.sam_write1(
		sf.fp,
		sf.hdr,
		br.b,
	)
	if ret < 0 {
		return int(ret), fmt.Errorf("boom: failed to write record")
	}
	return int(ret), nil
}

// A bamIndex wraps a hts_idx_t.
type bamIndex struct {
	idx *C.hts_idx_t
}

// bamIndexBuild builds a BAM index file, filename.bai, from a bam file, filename. It returns an
// integer value and any error that occurred.
func bamIndexBuild(filename string) (ret int, err error) {
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))

	r := C.sam_index_build(
		(*C.char)(unsafe.Pointer(fn)),
		0, // min_shift parameter, 0 for default
	)
	if r < 0 {
		return int(r), fmt.Errorf("boom: failed to build index")
	}

	return int(r), nil
}

// bamIndexLoad loads a BAM index, returning a *bamIndex and any error that occurred.
// The bamIndex is created setting a finaliser that destroys the contained hts_idx_t.
func bamIndexLoad(filename string) (bi *bamIndex, err error) {
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))

	ip := C.hts_idx_load(
		(*C.char)(unsafe.Pointer(fn)),
		C.int(1), // HTS_FMT_BAI
	)
	if ip == nil {
		return nil, fmt.Errorf("boom: failed to load index")
	}
	bi = &bamIndex{idx: ip}
	runtime.SetFinalizer(bi, (*bamIndex).bamIndexDestroy)

	return
}

// bamIndexDestroy destroys the contained hts_idx_t, first checking for nil pointers.
func (bi *bamIndex) bamIndexDestroy() (err error) {
	if bi.idx == nil {
		return valueIsNil
	}

	C.hts_idx_destroy(bi.idx)
	bi.idx = nil

	return
}

// A bamFetchFn is called on each bamRecord found by bamFetch. The return value is used to indicate
// the iteration is complete.
type bamFetchFn func(*bamRecord) bool

// bamFetch calls fn on all BAM records within the interval [beg, end) of the reference sequence
// identified by tid. Note that beg >= 0 || beg = 0.
func (sf *samFile) bamFetch(bi *bamIndex, tid, beg, end int, fn bamFetchFn) (ret int, err error) {
	if sf.fp == nil || sf.hdr == nil || bi.idx == nil {
		return 0, valueIsNil
	}

	iter := C.sam_itr_queryi(bi.idx, C.int(tid), C.hts_pos_t(beg), C.hts_pos_t(end))
	if iter == nil {
		return 0, fmt.Errorf("boom: failed to create iterator")
	}
	defer C.hts_itr_destroy(iter)

	var br *bamRecord
	for {
		br, err = newBamRecord(nil)
		if err != nil {
			return
		}
		ret = int(C.sam_itr_next(sf.fp, iter, br.b))
		if ret < 0 {
			break
		}
		if fn(br) {
			break
		}
	}

	return
}

// Type header defines types that can be passed to samOpen as a SAM header or header filename.
type header interface {
	header() // No-op for interface definition.
}

// A bamHeader wraps a sam_hdr_t.
type bamHeader struct {
	bh *C.sam_hdr_t
}

// bamGetTid return the target id for for a reference sequence target matching the string, name.
func (bh *bamHeader) bamGetTid(name string) int {
	if bh.bh == nil {
		panic(valueIsNil)
	}

	sn := C.CString(name)
	defer C.free(unsafe.Pointer(sn))

	tid := C.sam_hdr_name2tid(
		bh.bh,
		(*C.char)(unsafe.Pointer(sn)),
	)

	return int(tid)
}

// nTargets returns the number of reference sequence targets described in the BAM header.
func (bh *bamHeader) nTargets() int32 {
	if bh.bh != nil {
		return int32(bh.bh.n_targets)
	}
	panic(valueIsNil)
}

// targetNames returns a slice of strings containing the names of the reference sequence
// targets described in the BAM header.
func (bh *bamHeader) targetNames() (n []string) {
	if bh.bh != nil {
		n = make([]string, bh.bh.n_targets)
		l := int(bh.bh.n_targets)
		var nPtrs []*C.char
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&nPtrs))
		sh.Cap = l
		sh.Len = l
		sh.Data = uintptr(unsafe.Pointer(bh.bh.target_name))

		for i, p := range nPtrs {
			n[i] = C.GoString(p)
		}

		return
	}
	panic(valueIsNil)
}

// targetLengths returns a slice of uint32 containing the lengths of the reference sequence
// targets described in the BAM header.
func (bh *bamHeader) targetLengths() []uint32 {
	if bh.bh != nil {
		l := int(bh.bh.n_targets)
		var unsafeLengths []uint32
		sh := (*reflect.SliceHeader)(unsafe.Pointer(&unsafeLengths))
		sh.Cap = l
		sh.Len = l
		sh.Data = uintptr(unsafe.Pointer(bh.bh.target_len))

		return append([]uint32(nil), unsafeLengths...)
	}
	panic(valueIsNil)
}

// text returns a string containing the full unparsed BAM header.
func (bh *bamHeader) text() (t string) {
	if bh.bh != nil {
		// Use sam_hdr_str to get the header text
		str := C.sam_hdr_str(bh.bh)
		if str == nil {
			return ""
		}
		return C.GoString(str)
	}
	panic(valueIsNil)
}

// header is a no-op function required to allow *bamHeader to satisfy the header interface.
func (bh *bamHeader) header() {}

// stringHeader is a string representation of a filename of a SAM header file.
type stringHeader string

// header is a no-op function required to allow stringHeader to satisfy the header interface.
func (sh stringHeader) header() {}

// textHeader is a []byte representation of a filename of a SAM header file.
type textHeader []byte

// header is a no-op function required to allow textHeader to satisfy the header interface.
func (th textHeader) header() {}

const (
	Paired        Flags = paired        // The read is paired in sequencing, no matter whether it is mapped in a pair.
	ProperPair    Flags = properPair    // The read is mapped in a proper pair.
	Unmapped      Flags = unmapped      // The read itself is unmapped; conflictive with BAM_FPROPER_PAIR.
	MateUnmapped  Flags = mateUnmapped  // The mate is unmapped.
	Reverse       Flags = reverse       // The read is mapped to the reverse strand.
	MateReverse   Flags = mateReverse   // The mate is mapped to the reverse strand.
	Read1         Flags = read1         // This is read1.
	Read2         Flags = read2         // This is read2.
	Secondary     Flags = secondary     // Not primary alignment.
	QCFail        Flags = qCFail        // QC failure.
	Duplicate     Flags = duplicate     // Optical or PCR duplicate.
	Supplementary Flags = supplementary // Supplementary alignment, indicates alignment is part of a chimeric alignment.
)

const (
	paired        = C.BAM_FPAIRED
	properPair    = C.BAM_FPROPER_PAIR
	unmapped      = C.BAM_FUNMAP
	mateUnmapped  = C.BAM_FMUNMAP
	reverse       = C.BAM_FREVERSE
	mateReverse   = C.BAM_FMREVERSE
	read1         = C.BAM_FREAD1
	read2         = C.BAM_FREAD2
	secondary     = C.BAM_FSECONDARY
	qCFail        = C.BAM_FQCFAIL
	duplicate     = C.BAM_FDUP
	supplementary = C.BAM_FSUPPLEMENTARY
)

// A Flags represents a BAM record's alignment FLAG field.
type Flags uint32

// String representation of BAM alignment flags:
//  0x001 - p - Paired
//  0x002 - P - ProperPair
//  0x004 - u - Unmapped
//  0x008 - U - MateUnmapped
//  0x010 - r - Reverse
//  0x020 - R - MateReverse
//  0x040 - 1 - Read1
//  0x080 - 2 - Read2
//  0x100 - s - Secondary
//  0x200 - f - QCFail
//  0x400 - d - Duplicate
//  0x800 - S - Supplementary
//
// Note that flag bits are represented high order to the right.
func (f Flags) String() string {
	// If 0x01 is unset, no assumptions can be made about 0x02, 0x08, 0x20, 0x40 and 0x80
	const pairedMask = ProperPair | MateUnmapped | MateReverse | MateReverse | Read1 | Read2
	if f&1 == 0 {
		f &^= pairedMask
	}

	const flags = "pPuUrR12sfdS"

	b := make([]byte, len(flags))
	for i, c := range flags {
		if f&(1<<uint(i)) != 0 {
			b[i] = byte(c)
		} else {
			b[i] = '-'
		}
	}

	return string(b)
}
