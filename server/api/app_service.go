package api

import (
	v1 "go-license-management/server/api/v1"
)

type AppService struct {
	v1 *v1.V1AppService
}

func (svc *AppService) GetV1Svc() *v1.V1AppService {
	return svc.v1
}

func (svc *AppService) SetV1Svc(v1Svc *v1.V1AppService) {
	svc.v1 = v1Svc
}
