package dto

type CreateUserRequest struct {
	Name           string `json:"name" binding:"required"`
	ShortName      string `json:"short_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Age            int    `json:"age" binding:"required,gte=0"`
	Password       string `json:"password" binding:"required,min=8"`
	OrganizationID int    `json:"organization_id" binding:"required"`
}
