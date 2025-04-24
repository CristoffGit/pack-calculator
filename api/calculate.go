package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CristoffGit/pack-calculator/internal/calc"
)

type Req struct{ Items int `json:"items"` }
type Resp struct{ Result map[int]int `json:"result"` }

func loadPacks() []int {
	if env := os.Getenv("PACK_SIZES"); env != "" {
		parts := strings.Split(env, ",")
		out := make([]int, 0, len(parts))
		for _, p := range parts {
			if n, _ := strconv.Atoi(strings.TrimSpace(p)); n > 0 {
				out = append(out, n)
			}
		}
		return out
	}
	// fallback compile-time default
	return []int{250, 500, 1000, 2000, 5000}
}

// Vercel → exported function signature
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var q Req
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	out, err := calc.Calculate(q.Items, loadPacks())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Resp{Result: out})
}
