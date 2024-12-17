package collections

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 有序字典示例
// om := NewOrderedMap[string, int]()
// om.Set("a", 1)
// om.Set("b", 2)
// om.ForEach(func(key string, value int) bool {
//     fmt.Printf("%s: %d\n", key, value) // 按插入顺序输出
//     return true
// })

// OrderedMapItem 有序字典项
type OrderedMapItem[K comparable, V any] struct {
	Key   K
	Value V
}

// OrderedMap 有序字典实现
type OrderedMap[K comparable, V any] struct {
	items    map[K]V
	keys     []K
	keyOrder map[K]int
}

// NewOrderedMap 创建新的有序字典
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		items:    make(map[K]V),
		keys:     make([]K, 0),
		keyOrder: make(map[K]int),
	}
}

// Set 设置键值对
func (m *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := m.items[key]; !exists {
		m.keys = append(m.keys, key)
		m.keyOrder[key] = len(m.keys) - 1
	}
	m.items[key] = value
}

// Get 获取值
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := m.items[key]
	return value, exists
}

// Delete 删除键值对
func (m *OrderedMap[K, V]) Delete(key K) {
	if order, exists := m.keyOrder[key]; exists {
		// 从 keys 切片中删除
		m.keys = append(m.keys[:order], m.keys[order+1:]...)
		// 更新后续键的顺序
		for i := order; i < len(m.keys); i++ {
			m.keyOrder[m.keys[i]] = i
		}
		// 删除键值对和顺序信息
		delete(m.items, key)
		delete(m.keyOrder, key)
	}
}

// Has 检查键是否存在
func (m *OrderedMap[K, V]) Has(key K) bool {
	_, exists := m.items[key]
	return exists
}

// Len 获取字典长度
func (m *OrderedMap[K, V]) Len() int {
	return len(m.items)
}

// Clear 清空字典
func (m *OrderedMap[K, V]) Clear() {
	m.items = make(map[K]V)
	m.keys = make([]K, 0)
	m.keyOrder = make(map[K]int)
}

// Keys 获取所有键（按插入顺序）
func (m *OrderedMap[K, V]) Keys() []K {
	result := make([]K, len(m.keys))
	copy(result, m.keys)
	return result
}

// Values 获取所有值（按键的插入顺序）
func (m *OrderedMap[K, V]) Values() []V {
	values := make([]V, len(m.keys))
	for i, key := range m.keys {
		values[i] = m.items[key]
	}
	return values
}

// Items 获取所有键值对（按插入顺序）
func (m *OrderedMap[K, V]) Items() []OrderedMapItem[K, V] {
	items := make([]OrderedMapItem[K, V], len(m.keys))
	for i, key := range m.keys {
		items[i] = OrderedMapItem[K, V]{
			Key:   key,
			Value: m.items[key],
		}
	}
	return items
}

// ForEach 遍历所有键值对（按插入顺序）
func (m *OrderedMap[K, V]) ForEach(fn func(key K, value V) bool) {
	for _, key := range m.keys {
		if !fn(key, m.items[key]) {
			break
		}
	}
}

// String 实现 Stringer 接口
func (m *OrderedMap[K, V]) String() string {
	return fmt.Sprintf("%v", m.Items())
}

// GoString 实现 GoStringer 接口
func (m *OrderedMap[K, V]) GoString() string {
	return fmt.Sprintf("%#v", m.Items())
}

// MarshalJSON 实现 json.Marshaler 接口
func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Items())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (m *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	var items []OrderedMapItem[K, V]
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	m.Clear()
	for _, item := range items {
		m.Set(item.Key, item.Value)
	}
	return nil
}

// First 获取第一个键值对
func (m *OrderedMap[K, V]) First() (K, V, error) {
	var key K
	var value V
	if len(m.keys) == 0 {
		return key, value, errors.New("map is empty")
	}
	key = m.keys[0]
	return key, m.items[key], nil
}

// Last 获取最后一个键值对
func (m *OrderedMap[K, V]) Last() (K, V, error) {
	var key K
	var value V
	if len(m.keys) == 0 {
		return key, value, errors.New("map is empty")
	}
	key = m.keys[len(m.keys)-1]
	return key, m.items[key], nil
}
