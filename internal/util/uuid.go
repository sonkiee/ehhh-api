package util

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(s string) (pgtype.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, errors.New("invalid uuid")
	}

	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}, nil
}
