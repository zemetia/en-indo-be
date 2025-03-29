package entity

type Kabupaten struct {
	ID         uint     `gorm:"type:int;primary_key;"`
	Name       string   `gorm:"type:varchar(255);not null"`
	ProvinsiID uint     `gorm:"type:int;not null"`
	Provinsi   Provinsi `gorm:"foreignKey:ProvinsiID"`

	Timestamp
}
