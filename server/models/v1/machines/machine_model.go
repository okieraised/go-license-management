package machines

import (
	"context"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/infrastructure/models/machine_attribute"
	"go-license-management/internal/server/v1/machines/models"
	"go-license-management/internal/utils"
	"go.opentelemetry.io/otel/trace"
)

type MachineRegistrationRequest struct {
	machine_attribute.MachineAttributeModel
}

func (req *MachineRegistrationRequest) Validate() error {
	if req.Fingerprint == nil {
		return comerrors.ErrMachineFingerprintIsEmpty
	}

	if req.LicenseKey == nil {
		return comerrors.ErrMachineLicenseIsEmpty
	} else {
		if len(utils.DerefPointer(req.LicenseKey)) != 44 {
			return comerrors.ErrMachineLicenseIsInvalid
		}
	}
	return nil
}

func (req *MachineRegistrationRequest) ToMachineRegistrationInput(ctx context.Context, tracer trace.Tracer, machineURI machine_attribute.MachineCommonURI) *models.MachineRegistrationInput {
	return &models.MachineRegistrationInput{
		TracerCtx:             ctx,
		Tracer:                tracer,
		MachineCommonURI:      machineURI,
		MachineAttributeModel: req.MachineAttributeModel,
	}
}

type MachineUpdateRequest struct {
	machine_attribute.MachineAttributeModel
	LicenseKey *string `json:"license_key"`
}

func (req *MachineUpdateRequest) Validate() error {
	if req.LicenseKey != nil {
		if len(utils.DerefPointer(req.LicenseKey)) != 44 {
			return comerrors.ErrMachineLicenseIsInvalid
		}
	}
	return nil
}

func (req *MachineUpdateRequest) ToMachineUpdateInput(ctx context.Context, tracer trace.Tracer, machineURI machine_attribute.MachineCommonURI) *models.MachineUpdateInput {

	return &models.MachineUpdateInput{
		TracerCtx: ctx,
		Tracer:    tracer,
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

type MachineDeactivateRequest struct{}

func (req *MachineDeactivateRequest) Validate() error {
	return nil
}

type MachineHeartbeatRequest struct{}

func (req *MachineHeartbeatRequest) Validate() error {
	return nil
}

type MachineActionsRequest struct {
	machine_attribute.MachineCommonURI
}

func (req *MachineActionsRequest) Validate() error {
	if req.MachineAction == nil {
		return comerrors.ErrMachineActionIsEmpty
	}

	return req.MachineCommonURI.Validate()
}

func (req *MachineActionsRequest) ToMachineActionsInput(ctx context.Context, tracer trace.Tracer, query machine_attribute.MachineActionsQueryParam) *models.MachineActionsInput {
	return &models.MachineActionsInput{
		TracerCtx:                ctx,
		Tracer:                   tracer,
		MachineCommonURI:         req.MachineCommonURI,
		MachineActionsQueryParam: query,
	}
}
