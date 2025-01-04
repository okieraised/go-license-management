package tokens

import (
	"github.com/uptrace/bun"
	"go-license-management/server/api"
)

type TokenRepository struct {
	database *bun.DB
}

func NewTokenRepository(ds *api.DataSource) *TokenRepository {
	return &TokenRepository{
		database: ds.GetDatabase(),
	}
}
