package skiplists

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/pkg/errors"

	"glookbs.github.com/storage"
)

func init() {
	// implement the DataStorage once
	var once sync.Once
	once.Do(func() {
		storage.DataStorage = *storage.New(New())
	})
}

var (
	MaxNodes = 512
	MaxLevel = 128
)

var (
	ErrSkipListIsFull       = errors.New("skip list is full")
	ErrSkipListDataNotFound = errors.New("data was not found")
)

type Node struct {
	key  int
	data any
	next []*Node
}

type SkipList struct {
	head   *Node
	level  int
	length int
}

func newNode(key, level int, data any) *Node {
	return &Node{
		key:  key,
		data: data,
		next: make([]*Node, level),
	}
}

func New() *SkipList {
	head := newNode(-1, MaxLevel, nil)
	return &SkipList{head, 1, 0}
}

func (sl *SkipList) Count() int {
	return sl.len()
}

func (sl *SkipList) len() int {
	return sl.length
}

// randomLevel returns the level with probability of coin flips
func randomLevel() int {
	level := 1
	for rand.Float32() < 0.5 && level < MaxLevel {
		level++
	}
	return level
}

// Insert inserts data and returns id if success, otherwise id=-1 with error
func (sl *SkipList) Insert(data any) (int, error) {
	id := sl.len() + 1
	if err := sl.insert(id, data); err != nil {
		return -1, err
	}
	return id, nil
}

func (sl *SkipList) insert(key int, data any) error {
	if sl.length >= MaxNodes {
		return ErrSkipListIsFull
	}

	update := make([]*Node, MaxLevel)
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		update[i] = current
	}

	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.head
		}
		sl.level = level
	}

	node := newNode(key, level, data)
	for i := 0; i < level; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node
	}

	sl.length++
	return nil
}

func (sl *SkipList) search(key int) (any, error) {
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
	}

	if current.next[0] != nil && current.next[0].key == key {
		return current.next[0].data, nil
	}

	return nil, ErrSkipListDataNotFound
}

// display shows the structure(only for debug)
func (list *SkipList) display() {
	for i := list.level - 1; i >= 0; i-- {
		current := list.head.next[i]
		fmt.Printf("Level %d: ", i)
		for current != nil {
			fmt.Printf("(%d, %v) -> ", current.key, current.data)
			current = current.next[i]
		}
		fmt.Println("nil")
	}
}

func (sl *SkipList) Delete(key int) bool {
	update := make([]*Node, MaxLevel)
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		update[i] = current
	}

	target := current.next[0]
	if target != nil && target.key == key {
		for i := 0; i < sl.level; i++ {
			if update[i].next[i] != target {
				break
			}
			update[i].next[i] = target.next[i]
		}
		for sl.level > 1 && sl.head.next[sl.level-1] == nil {
			sl.level--
		}
		sl.length--
		return true
	}

	return false
}

// Range returns data with range, returns empty slice if no data
func (sl *SkipList) Range(i, j int) []any {
	j++
	result := []any{}
	current := sl.head
	start := (i - 1) * j
	end := start + j
	count := 0

	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < start {
			current = current.next[i]
		}
	}

	current = current.next[0]
	for current != nil && current.key < end && count < j {
		if current.key >= start {
			result = append(result, current.key)
			count++
		}
		current = current.next[0]
	}

	return result
}

// Update performs delete+insert
func (sl *SkipList) Update(id int, data any) error {
	if !sl.Delete(id) {
		return ErrSkipListDataNotFound
	}

	return sl.insert(id, data)
}
