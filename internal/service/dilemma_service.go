package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sonkiee/ehhh-api/internal/domain"
	"github.com/sonkiee/ehhh-api/internal/repository"
)

type DilemmaService struct {
	dRepo repository.DilemmaRepository
	vRepo repository.VoteRepository
}

func NewDilemmaService(dRepo repository.DilemmaRepository, vRepo repository.VoteRepository) *DilemmaService {
	return &DilemmaService{
		dRepo: dRepo,
		vRepo: vRepo,
	}
}

func (s *DilemmaService) Create(ctx context.Context, question, optionA, optionB string) (domain.Dilemma, error) {
	q := strings.TrimSpace(question)
	a := strings.TrimSpace(optionA)
	b := strings.TrimSpace(optionB)

	if len(q) < 5 {
		return domain.Dilemma{}, errors.New("question too short")
	}
	if a == "" || b == "" {
		return domain.Dilemma{}, errors.New("both options are required")
	}
	if strings.EqualFold(a, b) {
		return domain.Dilemma{}, errors.New(("options must be different"))
	}

	d := domain.Dilemma{
		ID:        uuid.NewString(),
		Question:  q,
		OptionA:   a,
		OptionB:   b,
		CreatedAt: time.Now().UTC(),
	}
	return s.dRepo.Create(ctx, d)
}

func (s *DilemmaService) List(ctx context.Context, limit, offset int) ([]domain.Dilemma, error) {
	return s.dRepo.List(ctx, limit, offset)
}

func (s *DilemmaService) Get(ctx context.Context, id string) (domain.Dilemma, int, int, error) {
	d, err := s.dRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Dilemma{}, 0, 0, err
	}
	a, b, err := s.vRepo.CountByDilemma(ctx, id)
	if err != nil {
		return domain.Dilemma{}, 0, 0, err
	}
	return d, a, b, nil
}

func (s *DilemmaService) Vote(ctx context.Context, id string, choice domain.VoteChoice) (int, int, error) {
	if choice != domain.VoteA && choice != domain.VoteB {
		return 0, 0, errors.New("inavlid vote choice")
	}

	// validate dilemma exist
	if _, err := s.dRepo.GetByID(ctx, id); err != nil {
		return 0, 0, err
	}

	v := domain.Vote{
		ID:        uuid.NewString(),
		DilemmaID: id,
		Choice:    choice,
		CreatedAt: time.Now().UTC(),
	}

	if _, err := s.vRepo.Create(ctx, v); err != nil {
		return 0, 0, err
	}
	return s.vRepo.CountByDilemma(ctx, id)
}
