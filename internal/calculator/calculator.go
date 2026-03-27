package calculator

import (
	"errors"
	"math"
	"sort"

	"github.com/linemk/pack-calculator/internal/store"
)

type Calculator struct {
	store *store.Store
}

func New(s *store.Store) *Calculator {
	return &Calculator{store: s}
}

func (c *Calculator) Calculate(order int) (map[int]int, error) {
	if order <= 0 {
		return nil, errors.New("order must be positive")
	}
	return solve(order, c.store.Get())
}

func solve(order int, sizes []int) (map[int]int, error) {
	if len(sizes) == 0 {
		return nil, errors.New("no pack sizes configured")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	maxPack := sizes[0]
	limit := order + maxPack

	dp := make([]int, limit+1)
	parent := make([]int, limit+1)
	for i := 1; i <= limit; i++ {
		dp[i] = math.MaxInt32
	}

	for i := 1; i <= limit; i++ {
		for _, p := range sizes {
			if i >= p && dp[i-p] != math.MaxInt32 && dp[i-p]+1 < dp[i] {
				dp[i] = dp[i-p] + 1
				parent[i] = p
			}
		}
	}

	for t := order; t <= limit; t++ {
		if dp[t] == math.MaxInt32 {
			continue
		}
		result := make(map[int]int)
		for cur := t; cur > 0; {
			p := parent[cur]
			result[p]++
			cur -= p
		}
		return result, nil
	}

	return nil, errors.New("no solution found")
}
