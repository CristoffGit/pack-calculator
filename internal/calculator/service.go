package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/CristoffGit/pack-calculator/internal/calc"
)

type Server struct {
	PackSizes []int
}

func New(packCfg []int) *Server { return &Server{PackSizes: packCfg} }

type req struct {
	Items int `json:"items"`
}
type resp struct {
	Result map[int]int `json:"result"`
}

func (s *Server) handleCalc(w http.ResponseWriter, r *http.Request) {
	var q req
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	out, err := calc.Calculate(q.Items, s.PackSizes)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp{Result: out})
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/calculate", s.handleCalc)
	mux.Handle("/", http.FileServer(http.Dir("./static"))) // UI later
	return mux
}

// Helper for main():
func LoadPacks() []int {
	if env := os.Getenv("PACK_SIZES"); env != "" {
		var ps []int
		_ = json.Unmarshal([]byte(env), &ps)
		return ps
	}
	b, _ := os.ReadFile("configs/pack.json")
	var ps []int
	_ = json.Unmarshal(b, &ps)
	return ps
}
