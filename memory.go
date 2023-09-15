package evm

import "github.com/holiman/uint256"

// evm的简单内存模型
type Memory struct {
	store       []byte
	lastGasCost uint64
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Set(offset, size uint64, val []byte) {
	if size > 0 {
		if offset+size > uint64(m.Len()) {
			panic("invalid memory: store empty")
		}
		copy(m.store[offset:offset+size], val)
	}
}

// solidity的int/int256/uint/uint256数值类型
func (m *Memory) Set32(offset uint64, val *uint256.Int) {
	if offset+32 > uint64(m.Len()) {
		panic("invalid memory: store empty")
	}
	b32 := val.Bytes32()
	copy(m.store[offset:], b32[:])
}

func (m *Memory) Resize(size uint64) {
	if uint64(m.Len()) < size {
		m.store = append(m.store, make([]byte, size-uint64(m.Len()))...)
	}
}

func (m *Memory) Len() int {
	return len(m.store)
}

func (m *Memory) Data() []byte {
	return m.store
}

// 复制src位置的数据到dst位置,数据可能会覆盖
// 默认slice容量足够大, 否则会panic
func (m *Memory) Copy(dst, src, len uint64) {
	if len == 0 {
		return
	}
	copy(m.store[dst:], m.store[src:src+len])
}

// 复制Memory中指定位置的数据到新的Slice中,并返回
func (m *Memory) GetCopy(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if m.Len() > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.store[offset:offset+size])
	}

	return
}

// 返回Memory指定位置数据的指针
func (m *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if m.Len() > int(offset) {
		return m.store[offset : offset+size]
	}

	return nil
}
