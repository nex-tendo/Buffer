package main

import (
	"fmt"
	"log"
)

type GoBuffer struct {
	data   []byte
	cursor int64
}

func NewGoBuffer(data []byte) *GoBuffer {
	return &GoBuffer{
		data:   data,
		cursor: 0,
	}
}

func (b *GoBuffer) ReadBit(out *byte, offset int64) error {
	byteIndex := offset / 8
	bitIndex := 7 - (offset % 8)

	if byteIndex >= int64(len(b.data)) {
		return fmt.Errorf("out of bounds")
	}

	*out = (b.data[byteIndex] >> uint(bitIndex)) & 1
	return nil
}

func (b *GoBuffer) ReadBits(out *uint64, off, n int64) error {
	var result uint64
	var bout byte

	for i := int64(0); i < n; i++ {
		err := b.ReadBit(&bout, off+i)
		if err != nil {
			return err
		}

		result = (result << 1) | uint64(bout)
	}
  
	*out = result
	return nil
}
