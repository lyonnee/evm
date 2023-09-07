package evm

import (
	"sync"

	"github.com/holiman/uint256"
)

// 使用sync.Pool(对象池)提高性能
var stackPool = sync.Pool{
	New: func() any {
		return &Stack{
			data: make([]uint256.Int, 16),
		}
	},
}

func newstack() *Stack {
	return stackPool.Get().(*Stack)
}

func returnStack(s *Stack) {
	s.data = s.data[:0]
	stackPool.Put(s)
}

// stack本身只是个FILO模型的Slice
type Stack struct {
	data []uint256.Int
}

func (s *Stack) Data() []uint256.Int {
	return s.data
}

// 添加数据到Stack顶部
func (s *Stack) push(d *uint256.Int) {
	s.data = append(s.data, *d)
}

// 取出Stack顶部的数据
func (s *Stack) pop() (ret uint256.Int) {
	// 取Slice最后一个值做返回值
	ret = s.data[len(s.data)-1]
	// 并移除
	s.data = s.data[:len(s.data)-1]

	return
}

func (s *Stack) len() int {
	return len(s.data)
}

// 交换指定位置和Stack顶部的值
func (s *Stack) swap(n int) {
	s.data[s.len()-n], s.data[s.len()-1] = s.data[s.len()-1], s.data[s.len()-n]
}

// 复制Stack指定位置的数据到顶部
func (s *Stack) dup(n int) {
	s.push(&s.data[s.len()-n])
}

// 返回但不删除Stack顶部的值
func (s *Stack) peek() *uint256.Int {
	return &s.data[s.len()-1]
}

// 返回但不删除指定位置的值
func (s *Stack) Back(n int) *uint256.Int {
	return &s.data[s.len()-1]
}
