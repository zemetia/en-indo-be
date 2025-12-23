package dto

type LaguRequest struct {
	Judul          string `json:"Judul" binding:"required"`
	Artis          string `json:"Artis"`
	YoutubeLink    string `json:"YoutubeLink"`
	Genre          string `json:"Genre"`
	Lirik          string `json:"Lirik" binding:"required"`
	Tags           string `json:"Tags"`
	NadaDasar      string `json:"NadaDasar"`
	TahunRilis     int    `json:"TahunRilis"`
	OriginalLaguID string `json:"OriginalLaguID"` // Handle as string to avoid UUID parsing error on empty string
}
