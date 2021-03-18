package main

import (
	"runtime"

	"github.com/klauspost/compress/zstd"
)

type compressor interface {
	Compress([]byte) []byte
	Decompress([]byte) ([]byte, error)
}

type zstdCompressor struct {
	dec *zstd.Decoder
	enc *zstd.Encoder
}

func (z zstdCompressor) Compress(src []byte) []byte {
	return z.enc.EncodeAll(src, nil)
}

func (z zstdCompressor) Decompress(src []byte) ([]byte, error) {
	return z.dec.DecodeAll(src, nil)
}

func newZstdCompressor() compressor {
	dec, _ := zstd.NewReader(nil, zstd.WithDecoderConcurrency(runtime.NumCPU()))
	enc, _ := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(runtime.NumCPU()))
	z := zstdCompressor{
		enc: enc,
		dec: dec,
	}
	return z
}
