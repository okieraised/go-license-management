package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/server/v1/policies/models"
	"go-license-management/internal/utils"
)

func (svc *PolicyService) updatePolicyField(ctx *gin.Context, input *models.PolicyUpdateInput, policy *entities.Policy) (*entities.Policy, error) {
	var err error

	// Generate new private/public key pair
	if input.Scheme != nil {
		scheme := utils.DerefPointer(input.Scheme)
		if policy.Scheme != scheme {
			var privateKey = ""
			var publicKey = ""
			svc.logger.GetLogger().Info(fmt.Sprintf("generating private/public key pair using [%s] algorithm", scheme))
			switch scheme {
			case constants.PolicySchemeED25519:
				privateKey, publicKey, err = utils.NewEd25519KeyPair()
				if err != nil {
					svc.logger.GetLogger().Error(err.Error())
					return policy, err
				}
			case constants.PolicySchemeRSA2048PKCS1:
				privateKey, publicKey, err = utils.NewRSA2048PKCS1KeyPair()
				if err != nil {
					svc.logger.GetLogger().Error(err.Error())
					return policy, err
				}
			default:
				svc.logger.GetLogger().Error(fmt.Sprintf("invalid supported sheme [%s]", scheme))
				return policy, err
			}
			policy.PrivateKey = privateKey
			policy.PublicKey = publicKey
		}
	}

	if input.Duration != nil {
		policy.Duration = utils.DerefPointer(input.Duration)
	}

	if input.MaxMachines != nil {
		policy.MaxMachines = utils.DerefPointer(input.MaxMachines)
	}

	if input.MaxUses != nil {
		policy.MaxUses = utils.DerefPointer(input.MaxUses)
	}

	if input.HeartbeatDuration != nil {
		policy.HeartbeatDuration = utils.DerefPointer(input.HeartbeatDuration)
	}

	if input.MaxUsers != nil {
		policy.MaxUsers = utils.DerefPointer(input.MaxUsers)
	}

	if input.ExpirationStrategy != nil {
		policy.ExpirationStrategy = utils.DerefPointer(input.ExpirationStrategy)
	}

	if input.AuthenticationStrategy != nil {
		policy.AuthenticationStrategy = utils.DerefPointer(input.AuthenticationStrategy)
	}

	if input.ExpirationBasis != nil {
		policy.ExpirationBasis = utils.DerefPointer(input.ExpirationBasis)
	}

	if input.OverageStrategy != nil {
		policy.OverageStrategy = utils.DerefPointer(input.OverageStrategy)
	}

	if input.RenewalBasis != nil {
		policy.RenewalBasis = utils.DerefPointer(input.RenewalBasis)
	}

	if input.HeartbeatBasis != nil {
		policy.HeartbeatBasis = utils.DerefPointer(input.HeartbeatBasis)
	}

	if input.CheckInInterval != nil {
		policy.CheckInInterval = utils.DerefPointer(input.CheckInInterval)
	}

	if input.RequireCheckIn != nil {
		policy.RequireCheckIn = utils.DerefPointer(input.RequireCheckIn)
	}
	if input.RequireHeartbeat != nil {
		policy.RequireHeartbeat = utils.DerefPointer(input.RequireHeartbeat)
	}
	if input.UsePool != nil {
		policy.UsePool = utils.DerefPointer(input.UsePool)
	}
	if input.Protected != nil {
		policy.Protected = utils.DerefPointer(input.Protected)
	}
	if input.RateLimited != nil {
		policy.RateLimited = utils.DerefPointer(input.RateLimited)
	}
	if input.Encrypted != nil {
		policy.Encrypted = utils.DerefPointer(input.Encrypted)
	}

	return policy, nil
}
