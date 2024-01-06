package storage

import (
	"errors"
	"sync"
)

// Enginer is the inteface for the data low-level contronl
type Enginer interface {
	// Insert inserts data into engine, and returns the id of the data if success
	Insert(data any) (int, error)
	// Count returns the number of data
	Count() int
	// Range returns data with i and j
	Range(i, j int) []any
	// Delete deletes the data with id, return false if data was nonexist
	Delete(i int) bool
	// Update updates data if it does exist
	Update(id int, data any) error
}

// New returns storage with injecting the Enginer
func New(enginer Enginer) *Storage {
	return &Storage{
		engine: enginer,
	}
}

// Storage is an object for low-level data engine controling, including thread-safe and error handling
// TODO: instead of mutex by the data engine, if it exports the locker
type Storage struct {
	mu     sync.RWMutex
	engine Enginer
}

func (s *Storage) Insert(data any) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Insert(data)
}

func (s *Storage) Count(data any) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.engine.Count()
}

func (s *Storage) Range(i, j int) []any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.engine.Range(i, j)
}

func (s *Storage) Delete(i int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.engine.Delete(i) {
		return errors.New("data is not exsit")
	}
	return nil
}

func (s *Storage) Update(id int, data any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.engine.Update(id, data)
}

// DataStorage is the default Storage
var DataStorage Storage
