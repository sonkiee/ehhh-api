package repository

import (
	"context"

	"github.com/sonkiee/ehhh-api/internal/domain"
)

type VoteRepository interface {
	Create(ctx context.Context, v domain.Vote) (domain.Vote, error)
	CountByDilemma(ctx context.Context, dilemmaID string) (countA int, countB int, err error)
}
