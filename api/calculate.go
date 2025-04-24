// api/calculate.go
package handler   // ✔ any name works as long as Handler is exported

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CristoffGit/pack-calculator/internal/calc"
)

/* ---------- request / response DTOs ---------- */

type request struct {
	Items int `json:"items"`
}
type response struct {
	Result map[int]int `json:"result"`
}

/* ---------- runtime config ---------- */

func loadPackSizes() []int {
	if env := os.Getenv("PACK_SIZES"); env != "" {
		// env var format: "250,500,1000,2000,5000"
		parts := strings.Split(env, ",")
		out := make([]int, 0, len(parts))
		for _, p := range parts {
			if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil && n > 0 {
				out = append(out, n)
			}
		}
		return out
	}
	return []int{250, 500, 1000, 2000, 5000} // fallback default
}

/* ---------- Vercel entry-point ---------- */

// Handler is automatically detected by Vercel's Go runtime.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	result, err := calc.Calculate(req.Items, loadPackSizes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Result: result})
}
