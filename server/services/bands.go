package services

import (
	"context"

	"github.com/vague2k/blkhell/server/database"
)

const bandsCtxKey ctxKey = "bands"

type BandsService struct {
	db *database.Queries
}

func NewBandsService(db *database.Queries) *BandsService {
	return &BandsService{db: db}
}

func (s *BandsService) BandsFromContext(ctx context.Context) ([]database.Band, bool) {
	b, ok := ctx.Value(bandsCtxKey).([]database.Band)
	return b, ok
}
