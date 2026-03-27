package calculator_test

import (
	"reflect"
	"testing"

	"github.com/linemk/pack-calculator/internal/calculator"
	"github.com/linemk/pack-calculator/internal/store"
)

type testCase struct {
	name   string
	order  int
	sizes  []int
	expect map[int]int
}

var cases = []testCase{
	{"single item", 1, []int{250, 500, 1000, 2000, 5000}, map[int]int{250: 1}},
	{"exact pack", 250, []int{250, 500, 1000, 2000, 5000}, map[int]int{250: 1}},
	{"just over small pack", 251, []int{250, 500, 1000, 2000, 5000}, map[int]int{500: 1}},
	{"two packs", 501, []int{250, 500, 1000, 2000, 5000}, map[int]int{500: 1, 250: 1}},
	{"large order", 12001, []int{250, 500, 1000, 2000, 5000}, map[int]int{5000: 2, 2000: 1, 250: 1}},
	{"edge case greedy fails", 500000, []int{23, 31, 53}, map[int]int{53: 9429, 31: 7, 23: 2}},
}

func newCalc(t *testing.T, sizes []int) *calculator.Calculator {
	t.Helper()
	s := store.New()
	if err := s.Set(sizes); err != nil {
		t.Fatalf("store.Set: %v", err)
	}
	return calculator.New(s)
}

func TestCalculate(t *testing.T) {
	t.Parallel()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := newCalc(t, tc.sizes)
			got, err := c.Calculate(tc.order)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.expect) {
				t.Errorf("order=%d sizes=%v\n  got  %v\n  want %v", tc.order, tc.sizes, got, tc.expect)
			}
		})
	}
}

func TestCalculateInvalidOrder(t *testing.T) {
	t.Parallel()
	c := newCalc(t, []int{250, 500})
	_, err := c.Calculate(-1)
	if err == nil {
		t.Fatal("expected error for negative order")
	}
}

func BenchmarkCalculateEdgeCase(b *testing.B) {
	s := store.New()
	s.Set([]int{23, 31, 53})
	c := calculator.New(s)
	b.ResetTimer()
	for range b.N {
		c.Calculate(500000)
	}
}
