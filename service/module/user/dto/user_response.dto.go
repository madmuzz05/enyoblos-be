package dto

import "github.com/madmuzz05/be-enyoblos/service/module/organization/entity"

type GetUserResponse struct {
	ID             int                  `json:"id"`
	Name           string               `json:"name"`
	ShortName      string               `json:"short_name"`
	Email          string               `json:"email"`
	Age            int                  `json:"age"`
	OrganizationID int                  `json:"organization_id"`
	Organization   *entity.Organization `json:"organization"`
}
