package dto

type KabupatenResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ProvinsiID uint   `json:"provinsi_id"`
	Provinsi   string `json:"provinsi"`
}
