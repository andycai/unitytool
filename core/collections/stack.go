package collections

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 栈示例
// stack := NewStack[int]()
// stack.Push(1)
// stack.Push(2)
// value, _ := stack.Pop() // 返回 2

// Stack 栈实现
type Stack[T any] struct {
	items []T
}

// NewStack 创建新栈
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// Push 压栈
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop 出栈
func (s *Stack[T]) Pop() (item T, err error) {
	if len(s.items) == 0 {
		return item, errors.New("stack is empty")
	}
	item = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// Peek 查看栈顶元素
func (s *Stack[T]) Peek() (item T, err error) {
	if len(s.items) == 0 {
		return item, errors.New("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

// IsEmpty 检查栈是否为空
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size 获取栈大小
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// Clear 清空栈
func (s *Stack[T]) Clear() {
	s.items = make([]T, 0)
}

// ToSlice 转换为切片
func (s *Stack[T]) ToSlice() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result
}

// FromSlice 从切片创建栈
func (s *Stack[T]) FromSlice(items []T) {
	s.items = make([]T, len(items))
	copy(s.items, items)
}

// String 实现 Stringer 接口
func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.items)
}

// GoString 实现 GoStringer 接口
func (s *Stack[T]) GoString() string {
	return fmt.Sprintf("%#v", s.items)
}

// MarshalJSON 实现 json.Marshaler 接口
func (s *Stack[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.items)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (s *Stack[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.items)
}
