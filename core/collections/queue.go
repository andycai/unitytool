package collections

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Queue 泛型队列实现
type Queue[T any] struct {
	items []T
}

// NewQueue 创建新队列
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// Enqueue 入队
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue 出队
func (q *Queue[T]) Dequeue() (item T, err error) {
	if len(q.items) == 0 {
		return item, errors.New("queue is empty")
	}
	item = q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// IsEmpty 检查队列是否为空
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size 获取队列大小
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// Clear 清空队列
func (q *Queue[T]) Clear() {
	q.items = make([]T, 0)
}

// Peek 查看队首元素
func (q *Queue[T]) Peek() (item T, err error) {
	if len(q.items) == 0 {
		return item, errors.New("queue is empty")
	}
	return q.items[0], nil
}

// ToSlice 转换为切片
func (q *Queue[T]) ToSlice() []T {
	return q.items
}

// FromSlice 从切片创建队列
func (q *Queue[T]) FromSlice(items []T) {
	q.items = make([]T, len(items))
	copy(q.items, items)
}

// String 实现 Stringer 接口
func (q *Queue[T]) String() string {
	return fmt.Sprintf("%v", q.items)
}

// GoString 实现 GoStringer 接口
func (q *Queue[T]) GoString() string {
	return fmt.Sprintf("%#v", q.items)
}

// MarshalJSON 实现 json.Marshaler 接口
func (q *Queue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.items)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (q *Queue[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &q.items)
}
