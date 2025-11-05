package handler

import "github.com/madmuzz05/be-enyoblos/service/module/organization/usecase"

type OrganizationHandler struct {
	OrganizationUsecase usecase.IOrganizationUsecase
}

func InitOrganizationHandler(organizationUsecase usecase.IOrganizationUsecase) *OrganizationHandler {
	return &OrganizationHandler{
		OrganizationUsecase: organizationUsecase,
	}
}
