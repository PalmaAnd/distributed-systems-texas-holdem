package api

import (
	"encoding/json"
	"net/http"

	"github.com/texas-holdem/backend/internal/poker"
)

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		HoleCards     []string `json:"hole_cards"`
		CommunityCards []string `json:"community_cards"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if len(req.HoleCards) != 2 {
		respondError(w, http.StatusBadRequest, "need exactly 2 hole cards")
		return
	}
	if len(req.CommunityCards) != 5 {
		respondError(w, http.StatusBadRequest, "need exactly 5 community cards")
		return
	}

	hole, err := poker.ParseCards(req.HoleCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	community, err := poker.ParseCards(req.CommunityCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := poker.EvaluateBestHand(hole, community)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"best_hand":  cardsToStrings(result.BestHand),
		"rank":       int(result.Rank),
		"rank_name":  result.RankName,
	})
}

func (s *Server) handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Hand1 struct {
			HoleCards     []string `json:"hole_cards"`
			CommunityCards []string `json:"community_cards"`
		} `json:"hand1"`
		Hand2 struct {
			HoleCards     []string `json:"hole_cards"`
			CommunityCards []string `json:"community_cards"`
		} `json:"hand2"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	hole1, err := poker.ParseCards(req.Hand1.HoleCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, "hand1: "+err.Error())
		return
	}
	comm1, err := poker.ParseCards(req.Hand1.CommunityCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, "hand1 community: "+err.Error())
		return
	}
	hole2, err := poker.ParseCards(req.Hand2.HoleCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, "hand2: "+err.Error())
		return
	}
	comm2, err := poker.ParseCards(req.Hand2.CommunityCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, "hand2 community: "+err.Error())
		return
	}
	if len(hole1) != 2 || len(comm1) != 5 || len(hole2) != 2 || len(comm2) != 5 {
		respondError(w, http.StatusBadRequest, "each hand needs 2 hole + 5 community cards")
		return
	}

	winner, h1, h2, err := poker.CompareHands(hole1, comm1, hole2, comm2)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	winnerStr := "tie"
	if winner == 1 {
		winnerStr = "hand1"
	} else if winner == 2 {
		winnerStr = "hand2"
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"winner":   winnerStr,
		"hand1":    map[string]any{"best_hand": cardsToStrings(h1.BestHand), "rank_name": h1.RankName},
		"hand2":    map[string]any{"best_hand": cardsToStrings(h2.BestHand), "rank_name": h2.RankName},
	})
}

func (s *Server) handleProbability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		HoleCards      []string `json:"hole_cards"`
		CommunityCards []string `json:"community_cards"`
		NumPlayers     int      `json:"num_players"`
		NumSims        int      `json:"num_sims"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if len(req.HoleCards) != 2 {
		respondError(w, http.StatusBadRequest, "need exactly 2 hole cards")
		return
	}
	if req.NumPlayers == 0 {
		req.NumPlayers = 2
	}
	if req.NumSims == 0 {
		req.NumSims = 10000
	}

	hole, err := poker.ParseCards(req.HoleCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	community, err := poker.ParseCards(req.CommunityCards)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	prob, err := poker.WinProbability(hole, community, req.NumPlayers, req.NumSims)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"win_probability": prob,
		"num_sims":        req.NumSims,
		"num_players":     req.NumPlayers,
	})
}

func cardsToStrings(c []poker.Card) []string {
	s := make([]string, len(c))
	for i, card := range c {
		s[i] = card.String()
	}
	return s
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
