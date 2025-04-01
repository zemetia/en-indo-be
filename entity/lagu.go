package entity

import "github.com/google/uuid"

type Lagu struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Judul       string    `gorm:"type:varchar(255);not null"`
	Artis       string    `gorm:"type:varchar(255);"`
	YoutubeLink string    `gorm:"type:varchar(255);"`
	Genre       string    `gorm:"type:varchar(100);"` // Genre lagu
	Lirik       string    `gorm:"type:text;not null"` // Lirik lagu
	TagLagu     []TagLagu `gorm:"many2many:lagu_tag_lagu;"`
	NadaDasar   string    `gorm:"type:char(2);"`
	TahunRilis  int       `gorm:"type:int;"` // Tahun rilis lagu

	Timestamp
}
