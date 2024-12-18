package gobuffer

import (
	"fmt"
)

type GoBuffer struct {
	buf  []byte
	off  int64
	cap  int64
	boff int64
	bcap int64
}

func NewGoBuffer(slices ...[]byte) *GoBuffer {
	buf := &GoBuffer{
		buf:  []byte{},
		off:  0,
		boff: 0,
	}

	if len(slices) == 0 {
		buf.Refresh()
		return buf
	}

	if len(slices) == 1 {
		buf.buf = slices[0]
		buf.Refresh()
		return buf
	}

	for _, slice := range slices {
		buf.buf = append(buf.buf, slice...)
	}

	buf.Refresh()
	return buf
}

func (b *GoBuffer) ReadBit(out *byte, offset int64) error {
	byteIndex := offset / 8
	bitIndex := 7 - (offset % 8)

	if byteIndex >= int64(len(b.buf)) {
		return fmt.Errorf("out of bounds")
	}

	*out = (b.buf[byteIndex] >> uint(bitIndex)) & 1
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

func (b *GoBuffer) Refresh() {
	b.bcap = int64(len(b.buf))

	b.cap = b.bcap * 8
}

// Original code of Grow from https://github.com/habak67/gobuffer/blob/master/buffer.go#L151
func (b *GoBuffer) Grow(size int64) {
	if size < 0 {
		panic(fmt.Errorf("Invalid size: cannot be negative"))
	}

	if size <= b.bcap {
		b.buf = b.buf[0 : b.off+size]
		b.Refresh()
		return
	}

	newCapacity := b.bcap * 2
	if newCapacity < size {
		newCapacity = size
	}

	tmp := make([]byte, newCapacity)
	copy(tmp, b.buf)

	b.buf = tmp
	b.Refresh()
}
