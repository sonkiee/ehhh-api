package domain

import "time"

type Dilemma struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	OptionA   string    `json:"option_a"`
	OptionB   string    `json:"option_b"`
	CreatedAt time.Time `json:"created_at"`
}

type VoteChoice string

const (
	VoteA VoteChoice = "A"
	VoteB VoteChoice = "B"
)

type Vote struct {
	ID        string     `json:"id"`
	DilemmaID string     `json:"dilemmaId"`
	Choice    VoteChoice `json:"choice"`
	CreatedAt time.Time  `json:"created_at"`
}
