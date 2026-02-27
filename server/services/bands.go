package services

import (
	"github.com/vague2k/blkhell/server/database"
)

const bandsCtxKey ctxKey = "bands"

type BandsService struct {
	db *database.Queries
}

func NewBandsService(db *database.Queries) *BandsService {
	return &BandsService{db: db}
}
