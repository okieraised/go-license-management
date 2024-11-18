package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-license-management/internal/response"
	"go-license-management/internal/server/v1/machines/models"
	"go-license-management/internal/utils"
)

func (svc *MachineService) checkout(ctx *gin.Context, input *models.MachineActionsInput) (*models.MachineActionCheckoutOutput, error) {

	// query machine info
	machine, err := svc.repo.SelectMachineByPK(ctx, uuid.MustParse(utils.DerefPointer(input.MachineID)))
	if err != nil {
		return nil, err
	}

	// Query license info
	license, err := svc.repo.SelectLicenseByPK(ctx, machine.LicenseID)
	if err != nil {
		return nil, err
	}

	// Query policy
	policy, err := svc.repo.SelectPolicyByPK(ctx, license.PolicyID)
	if err != nil {
		return nil, err
	}

	fmt.Println("machine", machine)

	return nil, nil
}

func (svc *MachineService) pingHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	return nil, nil
}

func (svc *MachineService) resetHeartbeat(ctx *gin.Context, input *models.MachineActionsInput) (*response.BaseOutput, error) {
	return nil, nil
}
