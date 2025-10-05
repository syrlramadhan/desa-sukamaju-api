package dto

type KontakResponse struct {
	IdKontak         string `json:"id_kontak"`
	Email            string `json:"email"`
	Telepon          string `json:"telepon"`
	Facebook         string `json:"facebook"`
	Youtube          string `json:"youtube"`
	Instagram        string `json:"instagram"`
}