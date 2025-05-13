package core

import (
	"errors"
	"fmt"
	"io"
)

const MaxBufferSize = 10 * 1024 * 1024 // 10MB

type DirectBuffer struct {
	// The underlying byte slice
	data []byte
}

// 使用缓存池
// 作用：减少内存分配和回收的开销
func NewDirectBuffer() *DirectBuffer {
	buf := bytePool.Get().([]byte)
	return &DirectBuffer{data: buf[:0]}
}

func (db *DirectBuffer) Reset() {
	bytePool.Put(db.data[:0])
	db.data = nil
}

// 读数据
func (db *DirectBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	p := db.data
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	// 剩余空间
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			if nMax > MaxBufferSize {
				return n - nStart, errors.New("DirectBuffer too large")
			}
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:]) // 从p[n:]开始读取数据到r的缓存区中
		// nn 为读取到的字节数
		n += nn
		if err != nil {
			db.data = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil // 返回增加的字节数和nil
			}
			return n, err
		}
	}
}

func (db *DirectBuffer) Len() int {
	if db == nil {
		return 0
	}
	return len(db.data)
}

func (db *DirectBuffer) String() string {
	if db == nil || len(db.data) == 0 || db.data == nil {
		return "DirectBuffer: empty"
	}
	return fmt.Sprintln("DirectBuffer: %v", db.data[:len(db.data)])
}

func (db *DirectBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(db.data)
	return int64(n), err
}

func (db *DirectBuffer) Bytes() []byte {
	return db.data
}

func (db *DirectBuffer) Write(p []byte) (int64, error) {
	db.data = append(db.data, p...)
	return int64(len(p)), nil
}

// 下面可以添加writeByte、writeString等方法
