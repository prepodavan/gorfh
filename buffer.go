package gorfh

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	defaultLenSize int = 4
)

type ErrBuffWritten struct {
	LenBefore int
	LenAfter int
	LenSize int
	Written int
}

func (e ErrBuffWritten) Error() string {
	return fmt.Sprintf(
		"buffer error: len before is %d, len size is %d, was written %d, len after is %d",
		e.LenBefore,
		e.LenSize,
		e.Written,
		e.LenAfter,
	)
}

type bufferWithLens struct {
	bytes.Buffer
	lens map[int]int
}

func (bwl *bufferWithLens) PutLen(index, length int) {
	if bwl.lens == nil {
		bwl.lens = make(map[int]int)
	}
	bwl.lens[index] = length
}

func (bwl *bufferWithLens) BytesOrdered(order binary.ByteOrder) (p []byte) {
	p = bwl.Buffer.Bytes()
	for i, length := range bwl.lens {
		order.PutUint32(p[i:], uint32(length))
	}
	return
}

func (bwl *bufferWithLens) WriteWithLen(f func(dst *bufferWithLens) (n int, err error)) (n int, err error) {
	oldLen := bwl.Len()
	n, err = bwl.Write(make([]byte, defaultLenSize)) // reserve space for length
	if err != nil {
		return
	} else if n != defaultLenSize {
		err = ErrBuffWritten{
			LenBefore: oldLen,
			LenAfter:  bwl.Len(),
			LenSize:   defaultLenSize,
			Written:   n,
		}
		return
	}
	n, err = f(bwl)
	if err != nil {
		return
	} else if oldLen + defaultLenSize + n != bwl.Len() {
		err = ErrBuffWritten{
			LenSize: defaultLenSize,
			LenBefore: oldLen,
			LenAfter: bwl.Len(),
			Written: n,
		}
		return
	}
	bwl.PutLen(oldLen, n)
	return
}
