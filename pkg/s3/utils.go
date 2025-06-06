package s3

import (
	"bytes"
)

type ByteReader struct {
	*bytes.Reader
}

func NewByteReader(data []byte) *ByteReader {
	return &ByteReader{Reader: bytes.NewReader(data)}
}

func (r *ByteReader) Close() error {
	return nil
}
