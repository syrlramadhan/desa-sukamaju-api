package dto

type KontakRequest struct {
	Email            string `json:"email"`
	Telepon          string `json:"telepon"`
	Facebook         string `json:"facebook"`
	Youtube          string `json:"youtube"`
	Instagram        string `json:"instagram"`
}