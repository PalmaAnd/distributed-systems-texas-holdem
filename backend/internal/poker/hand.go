package poker

import (
	"sort"
)

// Hand rank values (higher = better)
const (
	HighCard RankType = iota + 1
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

type RankType int

func (r RankType) String() string {
	switch r {
	case RoyalFlush:
		return "Royal Flush"
	case StraightFlush:
		return "Straight Flush"
	case FourOfAKind:
		return "Four of a Kind"
	case FullHouse:
		return "Full House"
	case Flush:
		return "Flush"
	case Straight:
		return "Straight"
	case ThreeOfAKind:
		return "Three of a Kind"
	case TwoPair:
		return "Two Pair"
	case OnePair:
		return "One Pair"
	case HighCard:
		return "High Card"
	default:
		return "Unknown"
	}
}

// EvaluatedHand holds the best 5-card hand and its rank.
type EvaluatedHand struct {
	BestHand []Card   `json:"best_hand"`
	Rank     RankType `json:"rank"`
	RankName string   `json:"rank_name"`
}

// EvaluateBestHand returns the best 5-card hand from 2 hole + up to 5 community cards.
func EvaluateBestHand(hole []Card, community []Card) (EvaluatedHand, error) {
	if len(hole) != 2 {
		return EvaluatedHand{}, &InvalidInputError{Msg: "need exactly 2 hole cards"}
	}
	if len(community) > 5 {
		return EvaluatedHand{}, &InvalidInputError{Msg: "max 5 community cards"}
	}
	all := append(append([]Card{}, hole...), community...)
	best := selectBestFive(all)
	rank := rankHand(best)
	return EvaluatedHand{
		BestHand: best,
		Rank:     rank,
		RankName: rank.String(),
	}, nil
}

func selectBestFive(cards []Card) []Card {
	if len(cards) <= 5 {
		return copyAndSort(cards)
	}
	// All C(7,5) or C(n,5) combinations
	return bestCombination(cards, 5)
}

func copyAndSort(c []Card) []Card {
	out := make([]Card, len(c))
	copy(out, c)
	sort.Slice(out, func(i, j int) bool { return out[i].Rank > out[j].Rank })
	return out
}

func bestCombination(cards []Card, k int) []Card {
	var best []Card
	var bestScore int64 = -1
	combine(cards, 0, k, nil, func(sel []Card) {
		sorted := copyAndSort(sel)
		score := handScore(sorted)
		if score > bestScore {
			bestScore = score
			best = sorted
		}
	})
	return best
}

func combine(cards []Card, start, k int, curr []Card, fn func([]Card)) {
	if k == 0 {
		out := make([]Card, len(curr))
		copy(out, curr)
		fn(out)
		return
	}
	for i := start; i <= len(cards)-k; i++ {
		next := make([]Card, len(curr)+1)
		copy(next, curr)
		next[len(curr)] = cards[i]
		combine(cards, i+1, k-1, next, fn)
	}
}

// handScore returns a comparable int64 for the 5-card hand (higher = better).
func handScore(c []Card) int64 {
	r := rankHand(c)
	// Score = (rank_type << 40) | (card values...)
	var score int64 = int64(r) << 40
	for i, c := range c {
		score |= int64(c.Rank) << (8 * (4 - i))
	}
	return score
}

func rankHand(c []Card) RankType {
	if len(c) != 5 {
		return HighCard
	}
	byRank := make(map[int]int)
	bySuit := make(map[byte]int)
	for _, card := range c {
		byRank[card.Rank]++
		bySuit[card.Suit]++
	}

	flush := len(bySuit) == 1
	straightVal := straightHigh(byRank)

	if flush && straightVal == RankA {
		return RoyalFlush
	}
	if flush && straightVal > 0 {
		return StraightFlush
	}
	if four := hasCount(byRank, 4); four {
		return FourOfAKind
	}
	if three, two := hasCount(byRank, 3), hasCount(byRank, 2); three && two {
		return FullHouse
	}
	if flush {
		return Flush
	}
	if straightVal > 0 {
		return Straight
	}
	if hasCount(byRank, 3) {
		return ThreeOfAKind
	}
	pairs := countPairs(byRank)
	if pairs == 2 {
		return TwoPair
	}
	if pairs == 1 {
		return OnePair
	}
	return HighCard
}

func hasCount(byRank map[int]int, n int) bool {
	for _, c := range byRank {
		if c == n {
			return true
		}
	}
	return false
}

func countPairs(byRank map[int]int) int {
	n := 0
	for _, c := range byRank {
		if c == 2 {
			n++
		}
	}
	return n
}

// straightHigh returns the high card of a straight, or 0 if not a straight.
func straightHigh(byRank map[int]int) int {
	if len(byRank) != 5 {
		return 0
	}
	allRanks := make([]int, 0, 5)
	for r := range byRank {
		allRanks = append(allRanks, r)
	}
	sort.Slice(allRanks, func(i, j int) bool { return allRanks[i] > allRanks[j] })
	// A-5 wheel
	if allRanks[0] == RankA && allRanks[1] == Rank5 && allRanks[2] == Rank4 &&
		allRanks[3] == Rank3 && allRanks[4] == Rank2 {
		return Rank5
	}
	for i := 0; i < len(allRanks)-1; i++ {
		if allRanks[i]-allRanks[i+1] != 1 {
			return 0
		}
	}
	return allRanks[0]
}

