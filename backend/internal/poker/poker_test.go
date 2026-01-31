package poker

import (
	"testing"
)

func TestParseCard(t *testing.T) {
	tests := []struct {
		in   string
		suit byte
		rank int
	}{
		{"HA", 'H', RankA},
		{"S7", 'S', Rank7},
		{"CT", 'C', Rank10},
		{"ha", 'H', RankA},
	}
	for _, tt := range tests {
		c, err := ParseCard(tt.in)
		if err != nil {
			t.Fatalf("ParseCard(%q): %v", tt.in, err)
		}
		if c.Suit != tt.suit || c.Rank != tt.rank {
			t.Errorf("ParseCard(%q) = %+v, want suit=%c rank=%d", tt.in, c, tt.suit, tt.rank)
		}
	}
}

func TestEvaluate_RoyalFlush(t *testing.T) {
	hole, _ := ParseCards([]string{"HA", "HK"})
	community, _ := ParseCards([]string{"HQ", "HJ", "HT", "S2", "D3"})
	h, err := EvaluateBestHand(hole, community)
	if err != nil {
		t.Fatal(err)
	}
	if h.Rank != RoyalFlush {
		t.Errorf("got rank %v, want Royal Flush", h.RankName)
	}
}

func TestEvaluate_StraightFlush(t *testing.T) {
	hole, _ := ParseCards([]string{"H9", "H8"})
	community, _ := ParseCards([]string{"H7", "H6", "H5", "S2", "D3"})
	h, err := EvaluateBestHand(hole, community)
	if err != nil {
		t.Fatal(err)
	}
	if h.Rank != StraightFlush {
		t.Errorf("got rank %v, want Straight Flush", h.RankName)
	}
}

func TestEvaluate_FullHouse(t *testing.T) {
	hole, _ := ParseCards([]string{"HA", "HK"})
	community, _ := ParseCards([]string{"SA", "DA", "SK", "S2", "D3"})
	h, err := EvaluateBestHand(hole, community)
	if err != nil {
		t.Fatal(err)
	}
	if h.Rank != FullHouse {
		t.Errorf("got rank %v, want Full House", h.RankName)
	}
}

func TestCompareHands(t *testing.T) {
	hole1, _ := ParseCards([]string{"HA", "HK"})
	comm1, _ := ParseCards([]string{"HQ", "HJ", "HT", "S2", "D3"})
	hole2, _ := ParseCards([]string{"C2", "C3"})
	comm2, _ := ParseCards([]string{"C4", "C5", "C6", "S7", "D8"})
	winner, h1, h2, err := CompareHands(hole1, comm1, hole2, comm2)
	if err != nil {
		t.Fatal(err)
	}
	if winner != 1 {
		t.Errorf("winner = %d, want 1 (hand1 Royal Flush beats straight flush)", winner)
	}
	_ = h1
	_ = h2
}
