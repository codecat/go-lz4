package lz4

// #cgo CFLAGS: -O3
// #include "src/lz4.h"
// #include "src/lz4.c"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// StreamDecoder lets you decode
type StreamDecoder interface {
	SetDictionary(dictionary []byte)
	Decompress(in, out []byte) error
	Close()
}

type streamDecoder struct {
	p *C.LZ4_streamDecode_t
}

// SetDictionary sets the dictionary for the decoder to use.
func (sd *streamDecoder) SetDictionary(dictionary []byte) {
	C.LZ4_setStreamDecode(sd.p, p(dictionary), clen(dictionary))
}

// Decompress with a known output size on the given decoder.
// len(out) should be equal to the length of the uncompressed out.
func (sd *streamDecoder) Decompress(in, out []byte) error {
	if int(C.LZ4_decompress_safe_continue(sd.p, p(in), p(out), clen(in), clen(out))) < 0 {
		return errors.New("Malformed compression stream")
	}

	return nil
}

// Close the decoder stream.
func (sd *streamDecoder) Close() {
	C.LZ4_freeStreamDecode(sd.p)
}

// NewStreamDecoder creates a new LZ4 stream decoder.
func NewStreamDecoder() StreamDecoder {
	return &streamDecoder{C.LZ4_createStreamDecode()}
}

// p gets a char pointer to the first byte of a []byte slice
func p(in []byte) *C.char {
	if len(in) == 0 {
		return (*C.char)(unsafe.Pointer(nil))
	}
	return (*C.char)(unsafe.Pointer(&in[0]))
}

// clen gets the length of a []byte slice as a char *
func clen(s []byte) C.int {
	return C.int(len(s))
}

// Uncompress with a known output size. len(out) should be equal to
// the length of the uncompressed out.
func Uncompress(in, out []byte) error {
	if int(C.LZ4_decompress_safe(p(in), p(out), clen(in), clen(out))) < 0 {
		return errors.New("Malformed compression stream")
	}

	return nil
}

// CompressBound calculates the size of the output buffer needed by
// Compress. This is based on the following macro:
//
// #define LZ4_COMPRESSBOUND(isize)
//      ((unsigned int)(isize) > (unsigned int)LZ4_MAX_INPUT_SIZE ? 0 : (isize) + ((isize)/255) + 16)
func CompressBound(in []byte) int {
	return len(in) + ((len(in) / 255) + 16)
}

// Compress compresses in and puts the content in out. len(out)
// should have enough space for the compressed data (use CompressBound
// to calculate). Returns the number of bytes in the out slice.
func Compress(in, out []byte) (outSize int, err error) {
	outSize = int(C.LZ4_compress_limitedOutput(p(in), p(out), clen(in), clen(out)))
	if outSize == 0 {
		err = fmt.Errorf("insufficient space for compression")
	}
	return
}
