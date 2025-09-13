package entity

type Provinsi struct {
	ID   uint   `gorm:"type:int;primary_key;"`
	Name string `gorm:"type:varchar(255);not null"`

	Timestamp
}
