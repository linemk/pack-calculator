package store

import (
	"errors"
	"sort"
	"sync"
)

var defaultSizes = []int{250, 500, 1000, 2000, 5000}

type Store struct {
	mu    sync.RWMutex
	sizes []int
}

func New() *Store {
	s := make([]int, len(defaultSizes))
	copy(s, defaultSizes)
	return &Store{sizes: s}
}

func (s *Store) Get() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]int, len(s.sizes))
	copy(out, s.sizes)
	return out
}

func (s *Store) Set(sizes []int) error {
	if len(sizes) == 0 {
		return errors.New("at least one pack size required")
	}
	seen := make(map[int]struct{}, len(sizes))
	for _, v := range sizes {
		if v <= 0 {
			return errors.New("pack sizes must be positive")
		}
		if _, ok := seen[v]; ok {
			return errors.New("duplicate pack size")
		}
		seen[v] = struct{}{}
	}
	cp := make([]int, len(sizes))
	copy(cp, sizes)
	sort.Sort(sort.Reverse(sort.IntSlice(cp)))

	s.mu.Lock()
	s.sizes = cp
	s.mu.Unlock()
	return nil
}
