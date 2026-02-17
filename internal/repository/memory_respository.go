package repository

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/sonkiee/ehhh-api/internal/domain"
)

var ErrNotFound = errors.New("not found")

type InMemoryDilemmaRepo struct {
	mu   sync.RWMutex
	byID map[string]domain.Dilemma
	list []domain.Dilemma // kept for ordering
}

func NewInMemoryDilemmaRepo() *InMemoryDilemmaRepo {
	return &InMemoryDilemmaRepo{
		byID: make(map[string]domain.Dilemma),
		list: make([]domain.Dilemma, 0),
	}
}

func (r *InMemoryDilemmaRepo) Create(ctx context.Context, d domain.Dilemma) (domain.Dilemma, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[d.ID] = d
	r.list = append(r.list, d)

	// keep newest first
	sort.SliceStable(r.list, func(i, j int) bool {
		return r.list[i].CreatedAt.After(r.list[j].CreatedAt)
	})

	return d, nil
}

func (r *InMemoryDilemmaRepo) GetByID(ctx context.Context, id string) (domain.Dilemma, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.byID[id]
	if !ok {
		return domain.Dilemma{}, ErrNotFound
	}
	return d, nil
}

func (r *InMemoryDilemmaRepo) List(ctx context.Context, limit, offset int) ([]domain.Dilemma, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if offset >= len(r.list) {
		return []domain.Dilemma{}, nil
	}

	end := offset + limit
	if end > len(r.list) {
		end = len(r.list)
	}

	out := make([]domain.Dilemma, end-offset)
	copy(out, r.list[offset:end])
	return out, nil
}
