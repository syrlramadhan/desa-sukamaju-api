package dto

type BeritaResponse struct {
	IdBerita           string   `json:"id_berita"`
	JudulBerita        string   `json:"judul_berita"`
	Kategori           string   `json:"kategori"`
	TanggalPelaksanaan string   `json:"tanggal_pelaksanaan"`
	Deskripsi          string   `json:"deskripsi"`
	GambarBerita       []string `json:"gambar_berita"`
	CreatedAt          string   `json:"created_at"`
}
