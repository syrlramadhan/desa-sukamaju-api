package dto

type AparatRequest struct {
	Nama           string `json:"nama"`
	Jabatan        string `json:"jabatan"`
	NoTelepon      string `json:"no_telepon"`
	Email          string `json:"email"`
	Status         string `json:"status"`
	PeriodeMulai   string `json:"periode_mulai"`
	PeriodeSelesai string `json:"periode_selesai"`
	Foto           string `json:"foto"`
}

type BulkDeleteAparatRequest struct {
	IDAparat []string `json:"id_aparat"`
}