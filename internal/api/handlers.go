package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/trucore-ai/x402fuel/internal/keystore"
	"github.com/trucore-ai/x402fuel/internal/policy"
	"github.com/trucore-ai/x402fuel/internal/types"
)

type Server struct {
	mu       sync.RWMutex
	wallets  map[string]*types.Wallet
	keystore *keystore.KeyStore
	policy   *policy.Engine
}

func NewServer(ks *keystore.KeyStore, pol *policy.Engine) *Server {
	return &Server{
		wallets:  make(map[string]*types.Wallet),
		keystore: ks,
		policy:   pol,
	}
}

func (s *Server) Register(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/wallets", s.createWallet)
		r.Get("/wallets", s.listWallets)
		r.Route("/wallets/{id}", func(r chi.Router) {
			r.Get("/", s.getWallet)
			r.Post("/pause", s.pauseWallet)
			r.Post("/resume", s.resumeWallet)
			r.Get("/transactions", s.getTransactions)
		})
		r.Get("/policy", s.getPolicy)
	})
}

func (s *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Label      string `json:"label"`
		Passphrase string `json:"passphrase"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	if req.Passphrase == "" {
		http.Error(w, `{"error":"passphrase required"}`, http.StatusBadRequest)
		return
	}

	addr, err := s.keystore.Generate(req.Passphrase)
	if err != nil {
		http.Error(w, `{"error":"failed to create wallet"}`, http.StatusInternalServerError)
		return
	}

	wallet := &types.Wallet{
		ID:      uuid.New().String(),
		Address: addr,
		Label:   req.Label,
	}

	s.mu.Lock()
	s.wallets[addr] = wallet
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      wallet.ID,
		"address": wallet.Address,
		"label":   wallet.Label,
	})
}

func (s *Server) listWallets(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addrs, err := s.keystore.List()
	if err != nil {
		http.Error(w, `{"error":"failed to list wallets"}`, http.StatusInternalServerError)
		return
	}

	var result []map[string]string
	for _, addr := range addrs {
		entry := map[string]string{"address": addr}
		if w, ok := s.wallets[addr]; ok {
			entry["id"] = w.ID
			entry["label"] = w.Label
		}
		result = append(result, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) getWallet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, wal := range s.wallets {
		if wal.ID == id || wal.Address == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"id":      wal.ID,
				"address": wal.Address,
				"label":   wal.Label,
			})
			return
		}
	}
	http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
}

func (s *Server) pauseWallet(w http.ResponseWriter, r *http.Request) {
	s.policy.SetPaused(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"paused": true,
		"message": "All payments paused. Use POST /resume to re-enable.",
	})
}

func (s *Server) resumeWallet(w http.ResponseWriter, r *http.Request) {
	s.policy.SetPaused(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"paused": false,
		"message": "Payments resumed.",
	})
}

func (s *Server) getTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]types.Transaction{})
}

func (s *Server) getPolicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.policy.GetPolicy())
}