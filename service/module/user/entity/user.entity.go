// User entity
package entity

type User struct {
	ID             int    `db:"id" json:"id"`
	Name           string `db:"name" json:"name" binding:"required"`
	ShortName      string `db:"short_name" json:"short_name" binding:"required"`
	Email          string `db:"email" json:"email" binding:"required,email"`
	Age            int    `db:"age" json:"age" binding:"required,min=0"`
	Password       string `db:"password" json:"password" binding:"required,min=8"`
	OrganizationID int    `db:"organization_id" json:"organization_id" binding:"required"`
}
