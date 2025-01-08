package models

import (
	"context"
	"github.com/google/uuid"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/models/machine_attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type MachineRegistrationInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	machine_attribute.MachineCommonURI
	machine_attribute.MachineAttributeModel
}

type MachineInfoOutput struct {
	ID              uuid.UUID              `json:"id"`
	LicenseKey      string                 `json:"license_key"`
	TenantName      string                 `json:"tenant_name"`
	Fingerprint     string                 `json:"fingerprint"`
	IP              string                 `json:"ip"`
	Hostname        string                 `json:"hostname"`
	Platform        string                 `json:"platform"`
	Name            string                 `json:"name"`
	Metadata        map[string]interface{} `json:"metadata"`
	Cores           int                    `json:"cores"`
	LastHeartbeatAt time.Time              `json:"last_heartbeat_at"`
	LastCheckOutAt  time.Time              `json:"last_check_out_at"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type MachineUpdateInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	machine_attribute.MachineCommonURI
	machine_attribute.MachineAttributeModel
}

type MachineRetrievalInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	machine_attribute.MachineCommonURI
}

type MachineDeleteInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	machine_attribute.MachineCommonURI
}

type MachineListInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	machine_attribute.MachineCommonURI
	constants.QueryCommonParam
}

type MachineListOutput struct {
	ID                   uuid.UUID              `json:"id"`
	LicenseID            uuid.UUID              `json:"license_id"`
	LicenseKey           string                 `json:"license_key"`
	TenantName           string                 `json:"tenant_name"`
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

type MachineActionsInput struct {
	TracerCtx context.Context
	Tracer    trace.Tracer
	machine_attribute.MachineCommonURI
	machine_attribute.MachineActionsQueryParam
}

type MachineActionCheckoutOutput struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Certificate string    `json:"certificate"`
	TTL         int       `json:"ttl"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
