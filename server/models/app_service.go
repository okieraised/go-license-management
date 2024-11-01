package models

import accountSvc "go-license-management/internal/server/v1/accounts/service"

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
	account *accountSvc.AccountService
}

func (v1 *V1AppService) GetAccount() *accountSvc.AccountService {
	return v1.account
}

func (v1 *V1AppService) SetAccount(svc *accountSvc.AccountService) {
	v1.account = svc
}
