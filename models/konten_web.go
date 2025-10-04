package models

type KontenWeb struct {
	IdKonten    string `json:"id_konten"`
	NamaWebsite string `json:"nama_website"`
	Logo        string `json:"logo"`
}

type Kontak struct {
	IdKontak         string `json:"id_kontak"`
	Email            string `json:"email"`
	Telepon          string `json:"telepon"`
	AlamatDesa       string `json:"alamat_desa"`
	LokasiMapsKantor string `json:"lokasi_maps_kantor"`
	Facebook         string `json:"facebook"`
	Youtube          string `json:"youtube"`
	Instagram        string `json:"instagram"`
}
