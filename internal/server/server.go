package server

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/CristoffGit/pack-calculator/internal/calc"
)

// ---------- server state -----------------------------------------------------

type Server struct {
	mu        sync.RWMutex // protects PackSizes at runtime
	PackSizes []int
}

func New(packs []int) *Server { return &Server{PackSizes: packs} }

// ---------- request / response shapes ---------------------------------------

type req struct{ Items int `json:"items"` }
type resp struct{ Result map[int]int `json:"result"` }

// ---------- /calculate -------------------------------------------------------

func (s *Server) handleCalc(w http.ResponseWriter, r *http.Request) {
	var q req
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	s.mu.RLock()
	out, err := calc.Calculate(q.Items, s.PackSizes)
	s.mu.RUnlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp{Result: out})
}

// ---------- /packs (GET & PUT) ----------------------------------------------

func (s *Server) handlePacks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet: // return current list
		s.mu.RLock()
		defer s.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.PackSizes)

	case http.MethodPut: // replace list with caller-supplied one
		var body struct {
			PackSizes []int `json:"packSizes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.PackSizes) == 0 {
			http.Error(w, "invalid packSizes", http.StatusBadRequest)
			return
		}
		s.mu.Lock()
		s.PackSizes = body.PackSizes
		s.mu.Unlock()
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// ---------- router -----------------------------------------------------------

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/calculate", s.handleCalc)
	mux.HandleFunc("/packs", s.handlePacks)               // NEW
	mux.Handle("/", http.FileServer(http.Dir("./static"))) // UI

	return mux
}

// ---------- helper to load packs at start-up --------------------------------

func LoadPacks() []int {
	// 1) env var overrides file
	if env := os.Getenv("PACK_SIZES"); env != "" {
		var ps []int
		_ = json.Unmarshal([]byte(env), &ps)
		return ps
	}
	// 2) fallback to configs/pack.json
	b, _ := os.ReadFile("configs/pack.json")
	var ps []int
	_ = json.Unmarshal(b, &ps)
	return ps
}
