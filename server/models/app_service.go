package models

import (
	accountSvc "go-license-management/internal/server/v1/accounts/service"
	authSvc "go-license-management/internal/server/v1/authentications/service"
	policySvc "go-license-management/internal/server/v1/policies/service"
	productSvc "go-license-management/internal/server/v1/products/service"
	tenantSvc "go-license-management/internal/server/v1/tenants/service"
)

type AppService struct {
	v1 *V1AppService
}

func (svc *AppService) GetV1Svc() *V1AppService {
	return svc.v1
}

func (svc *AppService) SetV1Svc(v1Svc *V1AppService) {
	svc.v1 = v1Svc
}

type V1AppService struct {
	account        *accountSvc.AccountService
	tenant         *tenantSvc.TenantService
	product        *productSvc.ProductService
	policy         *policySvc.PolicyService
	authentication *authSvc.AuthenticationService
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

func (v1 *V1AppService) GetAuth() *authSvc.AuthenticationService {
	return v1.authentication
}

func (v1 *V1AppService) SetAuth(svc *authSvc.AuthenticationService) {
	v1.authentication = svc
}
