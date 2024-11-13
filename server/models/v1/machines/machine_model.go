package machines

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/models/machine_attribute"
	"go-license-management/internal/server/v1/machines/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type MachineRegistrationRequest struct {
	machine_attribute.MachineAttributeModel
	LicenseID *string `json:"license_id"`
}

func (req *MachineRegistrationRequest) Validate() error {
	if req.Fingerprint == nil {
		return comerrors.ErrMachineFingerprintIsEmpty
	}

	if req.LicenseID == nil {
		return comerrors.ErrMachineLicenseIsEmpty
	} else {
		if _, err := uuid.Parse(utils.DerefPointer(req.LicenseID)); err != nil {
			return comerrors.ErrMachineLicenseIsInvalid
		}
	}
	return nil
}

func (req *MachineRegistrationRequest) ToMachineRegistrationInput(ctx context.Context, tracer trace.Tracer, machineURI machine_attribute.MachineCommonURI) *models.MachineRegistrationInput {
	return &models.MachineRegistrationInput{
		TracerCtx:             ctx,
		Tracer:                tracer,
		LicenseID:             req.LicenseID,
		MachineCommonURI:      machineURI,
		MachineAttributeModel: req.MachineAttributeModel,
	}
}

type MachineRetrievalRequest struct {
	machine_attribute.MachineCommonURI
}

func (req *MachineRetrievalRequest) Validate() error {
	if req.MachineID == nil {
		return comerrors.ErrMachineIDIsEmpty
	}
	return req.MachineCommonURI.Validate()
}

func (req *MachineRetrievalRequest) ToMachineRetrievalInput(ctx context.Context, tracer trace.Tracer) *models.MachineRetrievalInput {
	return &models.MachineRetrievalInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		MachineCommonURI: req.MachineCommonURI,
	}

}

type MachineDeletionRequest struct {
	machine_attribute.MachineCommonURI
}

func (req *MachineDeletionRequest) Validate() error {
	if req.MachineID == nil {
		return comerrors.ErrMachineIDIsEmpty
	}
	return req.MachineCommonURI.Validate()
}

func (req *MachineDeletionRequest) ToMachineDeletionInput(ctx context.Context, tracer trace.Tracer) *models.MachineDeleteInput {
	return &models.MachineDeleteInput{
		TracerCtx:        ctx,
		Tracer:           tracer,
		MachineCommonURI: req.MachineCommonURI,
	}

}

type MachineUpdateRequest struct{}

func (req *MachineUpdateRequest) Validate() error {
	return nil
}

type MachineDeactivateRequest struct{}

func (req *MachineDeactivateRequest) Validate() error {
	return nil
}

type MachineHeartbeatRequest struct{}

func (req *MachineHeartbeatRequest) Validate() error {
	return nil
}
