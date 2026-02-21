package entity

type Organization struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	ShortName string `db:"short_name" json:"short_name"`
	Address   string `db:"address" json:"address,omitempty"`
}

func (Organization) TableName() string {
	return "organizations"
}
