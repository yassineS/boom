// Copyright ©2025 The boom Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boom

/*
#cgo CFLAGS: -g -O2 -fPIC -m64 -pthread
#cgo LDFLAGS: -lhts
#include <htslib/vcf.h>
#include <htslib/hts.h>
*/
import "C"

import (
	"fmt"
	"io"
	"runtime"
	"unsafe"
)

// VCFFile represents an open VCF/BCF file.
type VCFFile struct {
	fp  *C.htsFile
	hdr *C.bcf_hdr_t
}

// OpenVCF opens a VCF or BCF file for reading.
func OpenVCF(filename string) (*VCFFile, error) {
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))

	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))

	fp := C.hts_open(fn, mode)
	if fp == nil {
		return nil, fmt.Errorf("boom: failed to open VCF file %s", filename)
	}

	hdr := C.bcf_hdr_read(fp)
	if hdr == nil {
		C.hts_close(fp)
		return nil, fmt.Errorf("boom: failed to read VCF header from %s", filename)
	}

	vf := &VCFFile{fp: fp, hdr: hdr}
	runtime.SetFinalizer(vf, (*VCFFile).Close)

	return vf, nil
}

// Close closes the VCF file and frees associated resources.
func (vf *VCFFile) Close() error {
	if vf.fp == nil {
		return valueIsNil
	}
	runtime.SetFinalizer(vf, nil)

	if vf.hdr != nil {
		C.bcf_hdr_destroy(vf.hdr)
		vf.hdr = nil
	}

	C.hts_close(vf.fp)
	vf.fp = nil

	return nil
}

// VCFRecord represents a single variant record from a VCF file.
type VCFRecord struct {
	rec *C.bcf1_t
}

// newVCFRecord creates a new VCF record.
func newVCFRecord() (*VCFRecord, error) {
	rec := C.bcf_init()
	if rec == nil {
		return nil, couldNotAllocate
	}

	vr := &VCFRecord{rec: rec}
	runtime.SetFinalizer(vr, (*VCFRecord).destroy)

	return vr, nil
}

// destroy frees the VCF record.
func (vr *VCFRecord) destroy() {
	if vr.rec != nil {
		C.bcf_destroy(vr.rec)
		vr.rec = nil
	}
}

// Read reads the next variant record from the VCF file.
func (vf *VCFFile) Read() (*VCFRecord, error) {
	if vf.fp == nil || vf.hdr == nil {
		return nil, valueIsNil
	}

	vr, err := newVCFRecord()
	if err != nil {
		return nil, err
	}

	ret := C.bcf_read(vf.fp, vf.hdr, vr.rec)
	if ret < 0 {
		return nil, io.EOF
	}

	// Unpack the record for easier access
	C.bcf_unpack(vr.rec, C.BCF_UN_ALL)

	return vr, nil
}

// Chrom returns the chromosome/contig name for the variant.
func (vr *VCFRecord) Chrom(hdr *C.bcf_hdr_t) string {
	if vr.rec == nil || hdr == nil {
		return ""
	}
	name := C.bcf_hdr_id2name(hdr, C.int(vr.rec.rid))
	if name == nil {
		return ""
	}
	return C.GoString(name)
}

// Pos returns the 0-based position of the variant.
func (vr *VCFRecord) Pos() int64 {
	if vr.rec == nil {
		return -1
	}
	return int64(vr.rec.pos)
}

// ID returns the variant ID.
func (vr *VCFRecord) ID() string {
	if vr.rec == nil {
		return ""
	}
	if vr.rec.d.id == nil {
		return ""
	}
	return C.GoString(vr.rec.d.id)
}

// Samples returns the number of samples in the VCF header.
func (vf *VCFFile) Samples() int {
	if vf.hdr == nil {
		return 0
	}
	return int(C.bcf_hdr_nsamples(vf.hdr))
}

// SampleNames returns the names of samples in the VCF file.
func (vf *VCFFile) SampleNames() []string {
	if vf.hdr == nil {
		return nil
	}

	n := int(C.bcf_hdr_nsamples(vf.hdr))
	if n == 0 {
		return nil
	}

	samples := make([]string, n)
	for i := 0; i < n; i++ {
		samples[i] = C.GoString(*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(vf.hdr.samples)) + uintptr(i)*unsafe.Sizeof((*C.char)(nil)))))
	}

	return samples
}
