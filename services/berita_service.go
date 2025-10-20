package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/models"
	"github.com/syrlramadhan/desa-sukamaju-api/repositories"
)

type BeritaService interface {
	CreateBerita(ctx context.Context, r *http.Request, beritaReq dto.BeritaRequest) (int, error)
	CreatePhoto(ctx context.Context, r *http.Request, photoReq dto.GaleriRequest) (int, error)
	GetAllBerita(ctx context.Context) ([]dto.BeritaResponse, int, error)
	GetBeritaById(ctx context.Context, idBerita string) (dto.BeritaResponse, int, error)
	UpdateBerita(ctx context.Context, idBerita string, beritaReq dto.BeritaRequest) (int, error)
	DeleteBerita(ctx context.Context, idBerita string) (int, error)
	DeletePhotoByFilename(ctx context.Context, filename string) (int, error)
	BulkDeletePhoto(ctx context.Context, filenames []string) (int, error)
}

type beritaServiceImpl struct {
	repo repositories.BeritaRepository
	DB   *sql.DB
}

func NewBeritaService(repo repositories.BeritaRepository, db *sql.DB) BeritaService {
	return &beritaServiceImpl{
		repo: repo,
		DB:   db,
	}
}

// CreateBerita implements BeritaService.
func (b *beritaServiceImpl) CreateBerita(ctx context.Context, r *http.Request, beritaReq dto.BeritaRequest) (int, error) {
	// ... (kode sebelumnya tetap sama)
	// Parse multipart form dengan buffer yang lebih besar (100MB)
	err := r.ParseMultipartForm(100 << 20) // 100MB
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to parse form: %v", err)
	}

	// Ambil data form
	beritaReq.JudulBerita = r.FormValue("judul_berita")
	beritaReq.Kategori = r.FormValue("kategori")
	beritaReq.TanggalPelaksanaan = r.FormValue("tanggal_pelaksanaan")
	beritaReq.Deskripsi = r.FormValue("deskripsi")

	// Validasi input
	if beritaReq.JudulBerita == "" {
		return http.StatusBadRequest, fmt.Errorf("judul berita tidak boleh kosong")
	}

	// Ambil multiple files dari form dengan pengecekan yang lebih safe
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return http.StatusBadRequest, fmt.Errorf("form data tidak valid")
	}

	fileHeaders, exists := r.MultipartForm.File["gambar_berita"]
	if !exists || len(fileHeaders) == 0 {
		return http.StatusBadRequest, fmt.Errorf("minimal satu gambar harus diupload")
	}

	// Batasi jumlah file maksimal (misal 10 files)
	if len(fileHeaders) > 10 {
		return http.StatusBadRequest, fmt.Errorf("maksimal 10 gambar dapat diupload sekaligus")
	}

	// Validasi dan proses multiple files
	var uploadedFilenames []string
	uploadDir := "uploads/berita"

	// Buat direktori jika belum ada
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal membuat direktori upload: %v", err)
	}

	beritaId := uuid.New().String()

	for i, fileHeader := range fileHeaders {
		// Validasi tipe file
		allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		isAllowed := false
		for _, allowedType := range allowedTypes {
			if ext == allowedType {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return http.StatusBadRequest, fmt.Errorf("tipe file tidak diizinkan untuk %s. Gunakan: jpg, jpeg, png, gif", fileHeader.Filename)
		}

		// Validasi ukuran file (maksimal 5MB per file untuk menghindari buffer issue)
		if fileHeader.Size > 5*1024*1024 {
			return http.StatusBadRequest, fmt.Errorf("ukuran file %s terlalu besar. Maksimal 5MB per file", fileHeader.Filename)
		}

		// Generate nama file unik
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("berita_%s_%d_%d%s", beritaId, timestamp, i+1, ext)

		// Buka file yang diupload
		file, err := fileHeader.Open()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("gagal membuka file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		// Path file tujuan
		filePath := filepath.Join(uploadDir, filename)

		// Buat file di server
		dst, err := os.Create(filePath)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("gagal membuat file %s: %v", filename, err)
		}
		defer dst.Close()

		// Copy file dengan limit untuk mencegah memory issues
		_, err = io.CopyN(dst, file, 5*1024*1024) // Limit 5MB
		if err != nil && err != io.EOF {
			return http.StatusInternalServerError, fmt.Errorf("gagal menyimpan file %s: %v", filename, err)
		}

		uploadedFilenames = append(uploadedFilenames, filename)
	}

	tx, err := b.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	berita := models.Berita{
		IdBerita:           beritaId,
		JudulBerita:        beritaReq.JudulBerita,
		Kategori:           beritaReq.Kategori,
		TanggalPelaksanaan: beritaReq.TanggalPelaksanaan,
		Deskripsi:          beritaReq.Deskripsi,
	}

	// Buat galeri untuk setiap gambar yang diupload
	var galeriList []models.Galeri
	for _, filename := range uploadedFilenames {
		galeri := models.Galeri{
			IdGaleri: uuid.New().String(),
			IdBerita: beritaId,
			Gambar:   filename,
		}
		galeriList = append(galeriList, galeri)
	}

	err = b.repo.AddBerita(ctx, tx, berita, galeriList)
	if err != nil {
		// Hapus file yang sudah diupload jika gagal menyimpan ke database
		for _, filename := range uploadedFilenames {
			os.Remove(filepath.Join(uploadDir, filename))
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal menambahkan berita: %v", err)
	}

	if err := tx.Commit(); err != nil {
		// Hapus file yang sudah diupload jika gagal commit
		for _, filename := range uploadedFilenames {
			os.Remove(filepath.Join(uploadDir, filename))
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}

// CreatePhoto implements BeritaService.
func (b *beritaServiceImpl) CreatePhoto(ctx context.Context, r *http.Request, photoReq dto.GaleriRequest) (int, error) {
	// ... (kode sebelumnya tetap sama)
	// Parse multipart form dengan buffer yang lebih besar (100MB)
	err := r.ParseMultipartForm(100 << 20) // 100MB
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to parse form: %v", err)
	}

	// Ambil data dari form
	idBerita := r.FormValue("id_berita")
	if idBerita == "" {
		return http.StatusBadRequest, fmt.Errorf("id_berita tidak boleh kosong")
	}

	// Ambil file dari form
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return http.StatusBadRequest, fmt.Errorf("form data tidak valid")
	}

	fileHeaders, exists := r.MultipartForm.File["gambar"]
	if !exists || len(fileHeaders) == 0 {
		return http.StatusBadRequest, fmt.Errorf("minimal satu gambar harus diupload")
	}

	// Batasi jumlah file maksimal (misal 5 files untuk AddPhoto)
	if len(fileHeaders) > 5 {
		return http.StatusBadRequest, fmt.Errorf("maksimal 5 gambar dapat diupload sekaligus untuk AddPhoto")
	}

	tx, err := b.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Validasi apakah berita exists
	berita, err := b.repo.GetBeritaById(ctx, tx, idBerita)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("berita dengan ID %s tidak ditemukan", idBerita)
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan berita: %v", err)
	}

	if berita.IdBerita == "" {
		return http.StatusNotFound, fmt.Errorf("berita dengan ID %s tidak ditemukan", idBerita)
	}

	// Setup upload directory
	uploadDir := "uploads/berita"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal membuat direktori upload: %v", err)
	}

	// Validasi dan proses multiple files
	var uploadedFilenames []string

	for i, fileHeader := range fileHeaders {
		// Validasi tipe file
		allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		isAllowed := false
		for _, allowedType := range allowedTypes {
			if ext == allowedType {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return http.StatusBadRequest, fmt.Errorf("tipe file tidak diizinkan untuk %s. Gunakan: jpg, jpeg, png, gif", fileHeader.Filename)
		}

		// Validasi ukuran file (maksimal 5MB per file)
		if fileHeader.Size > 5*1024*1024 {
			return http.StatusBadRequest, fmt.Errorf("ukuran file %s terlalu besar. Maksimal 5MB per file", fileHeader.Filename)
		}

		// Generate nama file unik
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("berita_%s_%d_additional_%d%s", idBerita, timestamp, i+1, ext)

		// Buka file yang diupload
		file, err := fileHeader.Open()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("gagal membuka file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		// Path file tujuan
		filePath := filepath.Join(uploadDir, filename)

		// Buat file di server
		dst, err := os.Create(filePath)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("gagal membuat file %s: %v", filename, err)
		}
		defer dst.Close()

		// Copy file dengan limit untuk mencegah memory issues
		_, err = io.CopyN(dst, file, 5*1024*1024) // Limit 5MB
		if err != nil && err != io.EOF {
			return http.StatusInternalServerError, fmt.Errorf("gagal menyimpan file %s: %v", filename, err)
		}

		uploadedFilenames = append(uploadedFilenames, filename)
	}

	// Buat galeri untuk setiap gambar yang diupload
	for _, filename := range uploadedFilenames {
		photo := models.Galeri{
			IdGaleri: uuid.New().String(),
			IdBerita: idBerita,
			Gambar:   filename,
		}

		// Simpan ke database
		err = b.repo.AddPhoto(ctx, tx, photo)
		if err != nil {
			// Hapus semua file yang sudah diupload jika gagal simpan ke database
			for _, uploadedFile := range uploadedFilenames {
				os.Remove(filepath.Join(uploadDir, uploadedFile))
			}
			return http.StatusInternalServerError, fmt.Errorf("gagal menambahkan foto: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		// Hapus semua file jika gagal commit
		for _, uploadedFile := range uploadedFilenames {
			os.Remove(filepath.Join(uploadDir, uploadedFile))
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}

// DeleteBerita implements BeritaService.
func (b *beritaServiceImpl) DeleteBerita(ctx context.Context, idBerita string) (int, error) {
	// Validasi input
	if idBerita == "" {
		return http.StatusBadRequest, fmt.Errorf("id berita tidak boleh kosong")
	}

	tx, err := b.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Ambil data foto berita untuk menghapus file
	photos, err := b.repo.GetPhotosByBeritaId(ctx, tx, idBerita)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan data foto: %v", err)
	}

	// Hapus berita dari database (termasuk galeri karena foreign key constraint)
	err = b.repo.DeleteBerita(ctx, tx, idBerita)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus berita: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	// Hapus file gambar dari storage setelah commit berhasil
	uploadDir := "uploads/berita"
	for _, photo := range photos {
		filePath := filepath.Join(uploadDir, photo.Gambar)
		if _, err := os.Stat(filePath); err == nil {
			os.Remove(filePath)
		}
	}

	return http.StatusOK, nil
}

// DeletePhotoByFilename implements BeritaService.
func (b *beritaServiceImpl) DeletePhotoByFilename(ctx context.Context, filename string) (int, error) {
	tx, err := b.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Hapus foto dari database
	err = b.repo.DeletePhotoByFilename(ctx, tx, filename)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus foto dari database: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	// Hapus file dari storage setelah commit berhasil
	uploadDir := "uploads/berita"
	filePath := filepath.Join(uploadDir, filename)
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	return http.StatusOK, nil
}

// BulkDeletePhoto implements BeritaService.
func (b *beritaServiceImpl) BulkDeletePhoto(ctx context.Context, filenames []string) (int, error) {
	// Validasi input
	if len(filenames) == 0 {
		return http.StatusBadRequest, fmt.Errorf("daftar nama file tidak boleh kosong")
	}

	tx, err := b.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Hapus multiple foto dari database
	err = b.repo.BulkDeletePhoto(ctx, tx, filenames)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus foto dari database: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	// Hapus file dari storage setelah commit berhasil
	uploadDir := "uploads/berita"
	for _, filename := range filenames {
		filePath := filepath.Join(uploadDir, filename)
		if _, err := os.Stat(filePath); err == nil {
			os.Remove(filePath)
		}
	}

	return http.StatusOK, nil
}

// GetAllBerita implements BeritaService.
func (b *beritaServiceImpl) GetAllBerita(ctx context.Context) ([]dto.BeritaResponse, int, error) {
	tx, err := b.DB.Begin()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Ambil semua berita dengan galeri menggunakan JOIN dari repository
	beritaList, err := b.repo.GetAllBerita(ctx, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan data berita: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	// Convert models.Berita ke dto.BeritaResponse
	var beritaResponseList []dto.BeritaResponse
	for _, berita := range beritaList {
		// Convert galeri ke format response
		var gambarBerita []string
		for _, galeri := range berita.GambarBerita {
			if galeri.Gambar != "" {
				// Tambahkan path lengkap untuk URL gambar
				gambarBerita = append(gambarBerita, galeri.Gambar)
			}
		}

		beritaResponse := dto.BeritaResponse{
			IdBerita:           berita.IdBerita,
			JudulBerita:        berita.JudulBerita,
			Kategori:           berita.Kategori,
			TanggalPelaksanaan: berita.TanggalPelaksanaan,
			Deskripsi:          berita.Deskripsi,
			GambarBerita:       gambarBerita,
			CreatedAt:          berita.TanggalPelaksanaan, // Menggunakan tanggal pelaksanaan sebagai created_at untuk sementara
		}

		beritaResponseList = append(beritaResponseList, beritaResponse)
	}

	return beritaResponseList, http.StatusOK, nil
}

// GetBeritaById implements BeritaService.
func (b *beritaServiceImpl) GetBeritaById(ctx context.Context, idBerita string) (dto.BeritaResponse, int, error) {
	// Validasi input
	if idBerita == "" {
		return dto.BeritaResponse{}, http.StatusBadRequest, fmt.Errorf("id berita tidak boleh kosong")
	}

	tx, err := b.DB.Begin()
	if err != nil {
		return dto.BeritaResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Ambil berita berdasarkan ID dari repository
	berita, err := b.repo.GetBeritaById(ctx, tx, idBerita)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.BeritaResponse{}, http.StatusNotFound, fmt.Errorf("berita dengan ID %s tidak ditemukan", idBerita)
		}
		return dto.BeritaResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan data berita: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return dto.BeritaResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	// Convert galeri ke format response
	var gambarBerita []string
	for _, galeri := range berita.GambarBerita {
		if galeri.Gambar != "" {
			// Tambahkan nama file gambar ke response
			gambarBerita = append(gambarBerita, galeri.Gambar)
		}
	}

	// Convert models.Berita ke dto.BeritaResponse
	beritaResponse := dto.BeritaResponse{
		IdBerita:           berita.IdBerita,
		JudulBerita:        berita.JudulBerita,
		Kategori:           berita.Kategori,
		TanggalPelaksanaan: berita.TanggalPelaksanaan,
		Deskripsi:          berita.Deskripsi,
		GambarBerita:       gambarBerita,
		CreatedAt:          berita.TanggalPelaksanaan, // Menggunakan tanggal pelaksanaan sebagai created_at untuk sementara
	}

	return beritaResponse, http.StatusOK, nil
}

// UpdateBerita implements BeritaService.
func (b *beritaServiceImpl) UpdateBerita(ctx context.Context, idBerita string, beritaReq dto.BeritaRequest) (int, error) {
	// Validasi input
	if idBerita == "" {
		return http.StatusBadRequest, fmt.Errorf("id berita tidak boleh kosong")
	}
	if beritaReq.JudulBerita == "" {
		return http.StatusBadRequest, fmt.Errorf("judul berita tidak boleh kosong")
	}
	if beritaReq.Kategori == "" {
		return http.StatusBadRequest, fmt.Errorf("kategori tidak boleh kosong")
	}
	if beritaReq.TanggalPelaksanaan == "" {
		return http.StatusBadRequest, fmt.Errorf("tanggal pelaksanaan tidak boleh kosong")
	}
	if beritaReq.Deskripsi == "" {
		return http.StatusBadRequest, fmt.Errorf("deskripsi tidak boleh kosong")
	}

	tx, err := b.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Validasi apakah berita exists
	_, err = b.repo.GetBeritaById(ctx, tx, idBerita)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("berita dengan ID %s tidak ditemukan", idBerita)
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan berita: %v", err)
	}

	// Update data berita
	berita := models.Berita{
		IdBerita:           idBerita,
		JudulBerita:        beritaReq.JudulBerita,
		Kategori:           beritaReq.Kategori,
		TanggalPelaksanaan: beritaReq.TanggalPelaksanaan,
		Deskripsi:          beritaReq.Deskripsi,
	}

	err = b.repo.UpdateBerita(ctx, tx, berita)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengupdate berita: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}