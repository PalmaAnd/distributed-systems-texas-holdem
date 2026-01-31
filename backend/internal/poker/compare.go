package poker

// CompareHands compares two hands (each: 2 hole + 5 community) and returns the winner.
// Returns 1 if hand1 wins, 2 if hand2 wins, 0 if tie.
func CompareHands(hole1, community1, hole2, community2 []Card) (int, EvaluatedHand, EvaluatedHand, error) {
	h1, err := EvaluateBestHand(hole1, community1)
	if err != nil {
		return -1, EvaluatedHand{}, EvaluatedHand{}, err
	}
	h2, err := EvaluateBestHand(hole2, community2)
	if err != nil {
		return -1, EvaluatedHand{}, EvaluatedHand{}, err
	}
	// Check for duplicate cards across both hands
	all1 := append(append([]Card{}, hole1...), community1...)
	all2 := append(append([]Card{}, hole2...), community2...)
	seen := make(map[string]bool)
	for _, c := range all1 {
		seen[c.String()] = true
	}
	for _, c := range all2 {
		if seen[c.String()] {
			return -1, EvaluatedHand{}, EvaluatedHand{}, &InvalidInputError{Msg: "cards cannot overlap between hands"}
		}
	}

	score1 := handScore(h1.BestHand)
	score2 := handScore(h2.BestHand)
	if score1 > score2 {
		return 1, h1, h2, nil
	}
	if score2 > score1 {
		return 2, h1, h2, nil
	}
	return 0, h1, h2, nil
}
