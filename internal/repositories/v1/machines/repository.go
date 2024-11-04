package machines

import (
	"github.com/uptrace/bun"
	"go-license-management/server/models"
)

type MachineRepository struct {
	database *bun.DB
}

func NewMachineRepository(ds *models.DataSource) *MachineRepository {
	return &MachineRepository{
		database: ds.GetDatabase(),
	}
}
