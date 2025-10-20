package dto

type BeritaRequest struct {
	JudulBerita        string   `json:"judul_berita"`
	Kategori           string   `json:"kategori"`
	TanggalPelaksanaan string   `json:"tanggal_pelaksanaan"`
	Deskripsi          string   `json:"deskripsi"`
	GambarBerita       []string `json:"gambar_berita"`
}
