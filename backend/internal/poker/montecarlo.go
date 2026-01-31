package poker

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Full deck of 52 cards as strings
var fullDeck []string

func init() {
	suits := []byte{'H', 'D', 'C', 'S'}
	ranks := []byte{'A', 'K', 'Q', 'J', 'T', '9', '8', '7', '6', '5', '4', '3', '2'}
	for _, s := range suits {
		for _, r := range ranks {
			fullDeck = append(fullDeck, string([]byte{s, r}))
		}
	}
}

// WinProbability runs Monte Carlo simulation and returns win probability for the given hand.
// hole: 2 hole cards, community: 0-5 known community cards, numPlayers: 2-10, numSims: simulations to run.
func WinProbability(hole []Card, community []Card, numPlayers, numSims int) (float64, error) {
	if len(hole) != 2 {
		return 0, &InvalidInputError{Msg: "need exactly 2 hole cards"}
	}
	if len(community) > 5 {
		return 0, &InvalidInputError{Msg: "max 5 community cards"}
	}
	if numPlayers < 2 || numPlayers > 10 {
		return 0, &InvalidInputError{Msg: "num_players must be 2-10"}
	}
	if numSims < 1 || numSims > 1000000 {
		return 0, &InvalidInputError{Msg: "num_sims must be 1-1000000"}
	}

	used := make(map[string]bool)
	for _, c := range hole {
		used[c.String()] = true
	}
	for _, c := range community {
		used[c.String()] = true
	}

	wins := 0
	for i := 0; i < numSims; i++ {
		deck := make([]string, 0, 52)
		for _, s := range fullDeck {
			if !used[s] {
				deck = append(deck, s)
			}
		}
		shuffle(deck)

		// Deal remaining community cards
		needed := 5 - len(community)
		commCards := append([]Card{}, community...)
		for j := 0; j < needed; j++ {
			c, _ := ParseCard(deck[j])
			commCards = append(commCards, c)
		}
		// Opponents get (numPlayers-1)*2 hole cards from deck[needed:]
		oppStart := needed
		oppCards := (numPlayers - 1) * 2
		if oppStart+oppCards > len(deck) {
			continue
		}

		// Evaluate our hand
		ourHand, _ := EvaluateBestHand(hole, commCards)
		ourScore := handScore(ourHand.BestHand)

		// Evaluate each opponent's hand
		weWin := true
		for k := 0; k < numPlayers-1; k++ {
			oh1, _ := ParseCard(deck[oppStart+k*2])
			oh2, _ := ParseCard(deck[oppStart+k*2+1])
			oppHand, _ := EvaluateBestHand([]Card{oh1, oh2}, commCards)
			oppScore := handScore(oppHand.BestHand)
			if oppScore > ourScore {
				weWin = false
				break
			}
			if oppScore == ourScore {
				// Tie: we split, count as half win for simplicity
				// Actually for "win probability" we typically mean strictly win. Let's not count ties as wins.
				weWin = false
				break
			}
		}
		if weWin {
			wins++
		}
	}
	return float64(wins) / float64(numSims), nil
}

func shuffle(a []string) {
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
