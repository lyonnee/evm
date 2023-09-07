package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryGet(t *testing.T) {
	odata := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	m := NewMemory()
	m.Resize(32)
	m.Set(0, 32, odata)

	bs1 := m.GetCopy(10, 5)
	bs1[0] = 0
	bs1[1] = 0
	bs1[2] = 0
	bs1[3] = 0
	bs1[4] = 0
	assert.Equal(t, odata, m.store)

	ndata := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	bs2 := m.GetPtr(10, 5)
	bs2[0] = 0
	bs2[1] = 0
	bs2[2] = 0
	bs2[3] = 0
	bs2[4] = 0
	assert.Equal(t, ndata, m.store)
}
