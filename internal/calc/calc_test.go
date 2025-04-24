package calc

import "testing"

func TestCalculate(t *testing.T) {
	packs := []int{250, 500, 1000, 2000, 5000}
	cases := []struct {
		items int
		want  map[int]int
	}{
		{1, map[int]int{250: 1}},
		{251, map[int]int{500: 1}},
		{501, map[int]int{500: 1, 250: 1}},
		{12001, map[int]int{5000: 2, 2000: 1, 250: 1}},
	}
	for _, c := range cases {
		got, _ := Calculate(c.items, packs)
		if !equal(got, c.want) {
			t.Errorf("items=%d got %v want %v", c.items, got, c.want)
		}
	}
}

func equal(a, b map[int]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
