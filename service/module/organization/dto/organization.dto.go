package dto

import "github.com/madmuzz05/be-enyoblos/service/module/organization/entity"

type PaginatedOrganizations struct {
	Data         []entity.Organization `json:"data"`
	TotalRecords int64                 `json:"total_records"`
}

// CreateOrganizationRequest - DTO untuk create organization
type CreateOrganizationRequest struct {
	Name      string `json:"name" binding:"required"`
	ShortName string `json:"short_name" binding:"required"`
	Address   string `json:"address"`
}

// UpdateOrganizationRequest - DTO untuk update organization
type UpdateOrganizationRequest struct {
	Name      string `json:"name" binding:"required"`
	ShortName string `json:"short_name" binding:"required"`
	Address   string `json:"address"`
}
