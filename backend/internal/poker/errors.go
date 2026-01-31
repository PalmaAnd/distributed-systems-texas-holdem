package poker

import "fmt"

type InvalidInputError struct {
	Msg string
}

func (e *InvalidInputError) Error() string {
	return e.Msg
}

func IsInvalidInput(err error) bool {
	_, ok := err.(*InvalidInputError)
	return ok
}

type DuplicateCardError struct {
	Card string
}

func (e *DuplicateCardError) Error() string {
	return fmt.Sprintf("duplicate card: %s", e.Card)
}
