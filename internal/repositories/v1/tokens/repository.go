package tokens

import (
	"github.com/uptrace/bun"
	"go-license-management/server/models"
)

type TokenRepository struct {
	database *bun.DB
}

func NewTokenRepository(ds *models.DataSource) *TokenRepository {
	return &TokenRepository{
		database: ds.GetDatabase(),
	}
}
