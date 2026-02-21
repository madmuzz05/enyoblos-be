package dto

type CreateUserRequest struct {
	Name           string `json:"name" validate:"required"`
	ShortName      string `json:"short_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Age            int    `json:"age" validate:"required,gte=0"`
	Password       string `json:"password" validate:"required,min=8"`
	OrganizationID int    `json:"organization_id" validate:"required"`
}
