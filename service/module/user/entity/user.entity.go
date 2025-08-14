// User entity
package entity

import orgEntity "github.com/madmuzz05/be-enyoblos/service/module/organization/entity"

type User struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
	ShortName      string `gorm:"type:varchar(100);not null" json:"short_name" binding:"required"`
	Email          string `gorm:"type:varchar(255);unique;not null" json:"email" binding:"required,email"`
	Age            int    `gorm:"not null" json:"age" binding:"required,min=0"`
	Password       string `gorm:"type:text;not null" json:"password" binding:"required,min=8"`
	OrganizationID *int   `gorm:"unique;not null" json:"organization_id" binding:"required"`

	// Relasi ke organization
	Organization *orgEntity.Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
}
