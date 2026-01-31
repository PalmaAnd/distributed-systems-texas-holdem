package poker

import (
	"fmt"
	"strings"
)

// Suit constants: H=Hearts, D=Diamonds, C=Clubs, S=Spades
const (
	SuitHearts   = 'H'
	SuitDiamonds = 'D'
	SuitClubs    = 'C'
	SuitSpades   = 'S'
)

// Rank values for comparison (2=2, ..., 14=Ace)
const (
	Rank2  = 2
	Rank3  = 3
	Rank4  = 4
	Rank5  = 5
	Rank6  = 6
	Rank7  = 7
	Rank8  = 8
	Rank9  = 9
	Rank10 = 10
	RankJ  = 11
	RankQ  = 12
	RankK  = 13
	RankA  = 14
)

var rankToChar = map[int]byte{
	Rank2: '2', Rank3: '3', Rank4: '4', Rank5: '5', Rank6: '6',
	Rank7: '7', Rank8: '8', Rank9: '9', Rank10: 'T',
	RankJ: 'J', RankQ: 'Q', RankK: 'K', RankA: 'A',
}

var charToRank = map[byte]int{
	'2': Rank2, '3': Rank3, '4': Rank4, '5': Rank5, '6': Rank6,
	'7': Rank7, '8': Rank8, '9': Rank9, 'T': Rank10,
	'J': RankJ, 'Q': RankQ, 'K': RankK, 'A': RankA,
}

// Card represents a playing card. Format: 2 chars e.g. "HA", "S7", "CT"
type Card struct {
	Suit byte // H, D, C, S
	Rank int  // 2-14 (Ace high)
}

// ParseCard parses a 2-char string (e.g. "HA", "S7") into a Card.
func ParseCard(s string) (Card, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if len(s) != 2 {
		return Card{}, fmt.Errorf("invalid card format: %q (need 2 chars)", s)
	}
	suit := s[0]
	if suit != 'H' && suit != 'D' && suit != 'C' && suit != 'S' {
		return Card{}, fmt.Errorf("invalid suit: %c (use H,D,C,S)", suit)
	}
	rank, ok := charToRank[s[1]]
	if !ok {
		return Card{}, fmt.Errorf("invalid rank: %c (use A,K,Q,J,T,9-2)", s[1])
	}
	return Card{Suit: suit, Rank: rank}, nil
}

// String returns the card as 2-char string (e.g. "HA").
func (c Card) String() string {
	return string([]byte{c.Suit, rankToChar[c.Rank]})
}

// ParseCards parses multiple cards.
func ParseCards(strs []string) ([]Card, error) {
	seen := make(map[string]bool)
	cards := make([]Card, 0, len(strs))
	for _, s := range strs {
		c, err := ParseCard(s)
		if err != nil {
			return nil, err
		}
		key := c.String()
		if seen[key] {
			return nil, fmt.Errorf("duplicate card: %s", key)
		}
		seen[key] = true
		cards = append(cards, c)
	}
	return cards, nil
}
