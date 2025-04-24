package calc

import "errors"

// ⬇⬇⬇  Moved OUTSIDE the function so every helper sees it
type node struct {
	items, packs int
	prevPack     int
}

// Calculate returns a map[packSize]quantity.
func Calculate(itemsOrdered int, packs []int) (map[int]int, error) {
	if itemsOrdered <= 0 {
		return nil, errors.New("items must be > 0")
	}

	max := itemsOrdered + maxOf(packs)
	dp  := make([]*node, max+1)
	dp[0] = &node{items: 0, packs: 0}

	for i := 1; i <= max; i++ {
		for _, p := range packs {
			if i-p < 0 || dp[i-p] == nil {
				continue
			}
			cand := node{
				items:    dp[i-p].items + p,
				packs:    dp[i-p].packs + 1,
				prevPack: p,
			}
			if better(&cand, dp[i], itemsOrdered) {
				cp := cand
				dp[i] = &cp
			}
		}
	}

	// find first feasible ≥ itemsOrdered
	idx := -1
	for i := itemsOrdered; i <= max && idx == -1; i++ {
		if dp[i] != nil {
			idx = i
		}
	}
	if idx == -1 {
		return nil, errors.New("no solution")
	}

	res := map[int]int{}
	for idx > 0 {
		p := dp[idx].prevPack
		res[p]++
		idx -= p
	}
	return res, nil
}

func better(a, b *node, want int) bool {
	if b == nil {
		return true
	}
	if a.items < b.items {
		return true
	}
	return a.items == b.items && a.packs < b.packs
}

func maxOf(xs []int) int {
	m := xs[0]
	for _, v := range xs[1:] {
		if v > m {
			m = v
		}
	}
	return m
}
