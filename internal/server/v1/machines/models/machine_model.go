package models

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/infrastructure/models/machine_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type MachineRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	machine_attribute.MachineCommonURI
	machine_attribute.MachineAttributeModel
	LicenseID *string `json:"license_id"`
}

type MachineRegistrationOutput struct {
	ID                   uuid.UUID              `json:"id"`
	TenantID             uuid.UUID              `json:"tenant_id"`
	LicenseID            uuid.UUID              `json:"license_id"`
	Fingerprint          string                 `json:"fingerprint"`
	IP                   string                 `json:"ip"`
	Hostname             string                 `json:"hostname"`
	Platform             string                 `json:"platform"`
	Name                 string                 `json:"name"`
	Metadata             map[string]interface{} `json:"metadata"`
	Cores                int                    `json:"cores"`
	LastHeartbeatAt      time.Time              `json:"last_heartbeat_at"`
	LastDeathEventSentAt time.Time              `json:"last_death_event_sent_at"`
	LastCheckOutAt       time.Time              `json:"last_check_out_at"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

type MachineUpdateInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
}

type MachineRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
}

type MachineDeleteInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
}

type MachineListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
}

type MachineActionsInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
}
