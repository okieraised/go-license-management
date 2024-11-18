package machine_attribute

import (
	"github.com/google/uuid"
	"go-license-management/internal/comerrors"
	"go-license-management/internal/constants"
	"go-license-management/internal/utils"
	"time"
)

type MachineCommonURI struct {
	TenantName    *string `uri:"tenant_name"`
	MachineID     *string `uri:"machine_id"`
	MachineAction *string `uri:"machine_action"`
}

func (req *MachineCommonURI) Validate() error {
	if req.TenantName == nil {
		return comerrors.ErrTenantNameIsEmpty
	}

	if req.MachineID != nil {
		if _, err := uuid.Parse(utils.DerefPointer(req.MachineID)); err != nil {
			return comerrors.ErrMachineIDIsInvalid
		}
	}

	if req.MachineAction != nil {
		if _, ok := constants.ValidMachineActionsMapper[utils.DerefPointer(req.MachineAction)]; !ok {
			return comerrors.ErrMachineActionIsInvalid
		}
	}

	return nil
}

type MachineActionsQueryParam struct {
	TTL *int `form:"ttl"`
}

func (req *MachineActionsQueryParam) Validate() error {
	if req.TTL == nil {
		req.TTL = utils.RefPointer(2592000)
	} else {
		ttl := utils.DerefPointer(req.TTL)
		if ttl < 3600 || ttl > 31556952 {
			return comerrors.ErrMachineActionCheckoutTTLIsInvalid
		}
	}

	return nil
}

// MachineAttributeModel contains information about the machine. Machines can be used to track and manage where your users are allowed to use your product.
type MachineAttributeModel struct {
	Fingerprint *string                `json:"fingerprint"` // The fingerprint of the machine. This can be an arbitrary string, but must be unique within the scope of the license it belongs to.
	Cores       *int                   `json:"cores"`       // The number of CPU cores for the machine.
	Name        *string                `json:"name"`        // The human-readable name of the machine.
	IP          *string                `json:"ip"`          // The IP of the machine.
	Hostname    *string                `json:"hostname"`    // The hostname of the machine.
	Platform    *string                `json:"platform"`    // The platform of the machine.
	Metadata    map[string]interface{} `json:"metadata"`    // Object containing machine metadata.
}

// MachineLicenseField contains information about the license
type MachineLicenseField struct {
	TenantID  string                 `json:"tenant_id"`
	ProductID string                 `json:"product_id"`
	PolicyID  string                 `json:"policy_id"`
	LicenseID string                 `json:"license_id"`
	MachineID string                 `json:"machine_id"`
	Metadata  map[string]interface{} `json:"metadata"`
	TTL       int                    `json:"ttl"`
	Expiry    time.Time              `json:"expiry"`
	CreatedAt time.Time              `json:"created_at"`
}

// MachineLicenseFileContent contains information about the license file
type MachineLicenseFileContent struct {
	Enc string `json:"enc"`
	Sig string `json:"sig"`
	Alg string `json:"alg"`
}
