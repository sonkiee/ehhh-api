package repository

import (
	"context"

	"github.com/sonkiee/ehhh-api/internal/domain"
)

type DilemmaRepository interface {
	Create(ctx context.Context, d domain.Dilemma) (domain.Dilemma, error)
	GetByID(ctx context.Context, id string) (domain.Dilemma, error)
	List(ctx context.Context, limit, offset int) ([]domain.Dilemma, error)
}
