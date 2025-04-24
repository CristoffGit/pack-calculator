// api/calculate.go
package main   // <- MUST be main

import (
	"encoding/json"
	"net/http"

	"github.com/CristoffGit/pack-calculator/internal/calc"
)

type request  struct{ Items int `json:"items"` }
type response struct{ Result map[int]int `json:"result"` }

func Handler(w http.ResponseWriter, r *http.Request) {  // <- exported
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var in request
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	out, err := calc.Calculate(in.Items, []int{250, 500, 1000, 2000, 5000})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Result: out})
}
