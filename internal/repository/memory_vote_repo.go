package repository

import (
	"context"
	"sync"

	"github.com/sonkiee/ehhh-api/internal/domain"
)

type InMemoryVoteRepo struct {
	mu sync.RWMutex

	// votesByDilemma keeps counts only (fast). If you later need per-user voting rules,
	// you'll store voter identity and enforce uniqueness.
	counts map[string]struct {
		A int
		B int
	}
}

func NewInMemoryVoteRepo() *InMemoryVoteRepo {
	return &InMemoryVoteRepo{
		counts: make(map[string]struct{ A, B int }),
	}
}

func (r *InMemoryVoteRepo) Create(ctx context.Context, v domain.Vote) (domain.Vote, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	c := r.counts[v.DilemmaID]
	if v.Choice == domain.VoteA {
		c.A++
	} else if v.Choice == domain.VoteB {
		c.B++
	}
	r.counts[v.DilemmaID] = c

	return v, nil
}

func (r *InMemoryVoteRepo) CountByDilemma(ctx context.Context, dilemmaID string) (int, int, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	c := r.counts[dilemmaID]
	return c.A, c.B, nil
}
