package service

import (
	"github.com/gin-gonic/gin"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/machines/models"
)

func (svc *MachineService) checkout(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {

	return nil, nil
}

func (svc *MachineService) pingHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *MachineService) resetHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	return nil, nil
}
