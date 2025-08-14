package entity

type Organization struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(255);unique;not null" json:"name"`
	ShortName string `gorm:"type:varchar(100);not null" json:"short_name"`
	Address   string `gorm:"type:text" json:"address,omitempty"`
}
