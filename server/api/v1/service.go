package v1

import (
	accountSvc "go-license-management/internal/services/v1/accounts/service"
	authSvc "go-license-management/internal/services/v1/authentications/service"
	entitlementSvc "go-license-management/internal/services/v1/entitlements/service"
	licenseSvc "go-license-management/internal/services/v1/licenses/service"
	machineSvc "go-license-management/internal/services/v1/machines/service"
	policySvc "go-license-management/internal/services/v1/policies/service"
	productSvc "go-license-management/internal/services/v1/products/service"
	tenantSvc "go-license-management/internal/services/v1/tenants/service"
)

type V1AppService struct {
	account        *accountSvc.AccountService
	tenant         *tenantSvc.TenantService
	product        *productSvc.ProductService
	policy         *policySvc.PolicyService
	entitlement    *entitlementSvc.EntitlementService
	machine        *machineSvc.MachineService
	authentication *authSvc.AuthenticationService
	license        *licenseSvc.LicenseService
}

func (v1 *V1AppService) GetAccount() *accountSvc.AccountService {
	return v1.account
}

func (v1 *V1AppService) SetAccount(svc *accountSvc.AccountService) {
	v1.account = svc
}

func (v1 *V1AppService) GetTenant() *tenantSvc.TenantService {
	return v1.tenant
}

func (v1 *V1AppService) SetTenant(svc *tenantSvc.TenantService) {
	v1.tenant = svc
}

func (v1 *V1AppService) GetProduct() *productSvc.ProductService {
	return v1.product
}

func (v1 *V1AppService) SetProduct(svc *productSvc.ProductService) {
	v1.product = svc
}

func (v1 *V1AppService) GetPolicy() *policySvc.PolicyService {
	return v1.policy
}

func (v1 *V1AppService) SetPolicy(svc *policySvc.PolicyService) {
	v1.policy = svc
}

func (v1 *V1AppService) GetEntitlement() *entitlementSvc.EntitlementService {
	return v1.entitlement
}

func (v1 *V1AppService) SetEntitlement(svc *entitlementSvc.EntitlementService) {
	v1.entitlement = svc
}

func (v1 *V1AppService) GetMachine() *machineSvc.MachineService {
	return v1.machine
}

func (v1 *V1AppService) SetMachine(svc *machineSvc.MachineService) {
	v1.machine = svc
}

func (v1 *V1AppService) GetAuth() *authSvc.AuthenticationService {
	return v1.authentication
}

func (v1 *V1AppService) SetAuth(svc *authSvc.AuthenticationService) {
	v1.authentication = svc
}

func (v1 *V1AppService) GetLicense() *licenseSvc.LicenseService {
	return v1.license
}

func (v1 *V1AppService) SetLicense(svc *licenseSvc.LicenseService) {
	v1.license = svc
}
