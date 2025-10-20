package repositories

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

type BeritaRepository interface {
	AddBerita(ctx context.Context, tx *sql.Tx, berita models.Berita, galeriList []models.Galeri) error
	AddPhoto(ctx context.Context, tx *sql.Tx, photo models.Galeri) error
	GetAllBerita(ctx context.Context, tx *sql.Tx) ([]models.Berita, error)
	GetBeritaById(ctx context.Context, tx *sql.Tx, idBerita string) (models.Berita, error)
	UpdateBerita(ctx context.Context, tx *sql.Tx, berita models.Berita) error
	DeleteBerita(ctx context.Context, tx *sql.Tx, idBerita string) error
}

type beritaRepositoryImpl struct{}

func NewBeritaRepository() BeritaRepository {
	return &beritaRepositoryImpl{}
}

// AddBerita implements BeritaRepository.
func (b *beritaRepositoryImpl) AddBerita(ctx context.Context, tx *sql.Tx, berita models.Berita, galeriList []models.Galeri) error {
	// Insert berita
	query := "INSERT INTO berita (id_berita, judul_berita, kategori, tanggal_pelaksanaan, deskripsi) VALUES (?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, berita.IdBerita, berita.JudulBerita, berita.Kategori, berita.TanggalPelaksanaan, berita.Deskripsi)
	if err != nil {
		return err
	}

	// Insert multiple galeri
	queryGaleri := "INSERT INTO galeri (id_galeri, id_berita, gambar) VALUES (?, ?, ?)"

	for _, galeri := range galeriList {
		_, err = tx.ExecContext(ctx, queryGaleri, galeri.IdGaleri, galeri.IdBerita, galeri.Gambar)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddPhoto implements BeritaRepository.
func (b *beritaRepositoryImpl) AddPhoto(ctx context.Context, tx *sql.Tx, photo models.Galeri) error {
	query := "INSERT INTO galeri (id_galeri, id_berita, gambar) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, photo.IdGaleri, photo.IdBerita, photo.Gambar)
	return err
}

// DeleteBerita implements BeritaRepository.
func (b *beritaRepositoryImpl) DeleteBerita(ctx context.Context, tx *sql.Tx, idBerita string) error {
	panic("unimplemented")
}

// GetAllBerita implements BeritaRepository.
func (b *beritaRepositoryImpl) GetAllBerita(ctx context.Context, tx *sql.Tx) ([]models.Berita, error) {
	query := `
		SELECT 
			b.id_berita, 
			b.judul_berita, 
			b.kategori, 
			b.tanggal_pelaksanaan, 
			b.deskripsi,
			COALESCE(g.id_galeri, '') as id_galeri,
			COALESCE(g.gambar, '') as gambar
		FROM berita b
		LEFT JOIN galeri g ON b.id_berita = g.id_berita
		ORDER BY b.id_berita, g.id_galeri
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk menyimpan berita berdasarkan ID
	beritaMap := make(map[string]*models.Berita)
	var beritaOrder []string // Untuk menjaga urutan

	for rows.Next() {
		var berita models.Berita
		var galeri models.Galeri

		err := rows.Scan(
			&berita.IdBerita,
			&berita.JudulBerita,
			&berita.Kategori,
			&berita.TanggalPelaksanaan,
			&berita.Deskripsi,
			&galeri.IdGaleri,
			&galeri.Gambar,
		)
		if err != nil {
			return nil, err
		}

		// Cek apakah berita sudah ada di map
		if existingBerita, exists := beritaMap[berita.IdBerita]; exists {
			// Jika berita sudah ada, tambahkan galeri ke list
			if galeri.IdGaleri != "" { // Hanya tambahkan jika ada galeri
				galeri.IdBerita = berita.IdBerita
				existingBerita.GambarBerita = append(existingBerita.GambarBerita, galeri)
			}
		} else {
			// Jika berita belum ada, buat entry baru
			newBerita := berita
			if galeri.IdGaleri != "" { // Hanya tambahkan jika ada galeri
				galeri.IdBerita = berita.IdBerita
				newBerita.GambarBerita = []models.Galeri{galeri}
			} else {
				newBerita.GambarBerita = []models.Galeri{} // Empty slice jika tidak ada galeri
			}

			beritaMap[berita.IdBerita] = &newBerita
			beritaOrder = append(beritaOrder, berita.IdBerita)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Convert map ke slice dengan urutan yang konsisten
	var beritaList []models.Berita
	for _, idBerita := range beritaOrder {
		beritaList = append(beritaList, *beritaMap[idBerita])
	}

	return beritaList, nil
}

// GetBeritaById implements BeritaRepository.
func (b *beritaRepositoryImpl) GetBeritaById(ctx context.Context, tx *sql.Tx, idBerita string) (models.Berita, error) {
	query := `
		SELECT 
			b.id_berita, 
			b.judul_berita, 
			b.kategori, 
			b.tanggal_pelaksanaan, 
			b.deskripsi,
			COALESCE(g.id_galeri, '') as id_galeri,
			COALESCE(g.gambar, '') as gambar
		FROM berita b
		LEFT JOIN galeri g ON b.id_berita = g.id_berita
		WHERE b.id_berita = ?
		ORDER BY g.id_galeri
	`

	rows, err := tx.QueryContext(ctx, query, idBerita)
	if err != nil {
		return models.Berita{}, err
	}
	defer rows.Close()

	var berita models.Berita
	var galeriList []models.Galeri
	isFirstRow := true

	for rows.Next() {
		var galeri models.Galeri

		err := rows.Scan(
			&berita.IdBerita,
			&berita.JudulBerita,
			&berita.Kategori,
			&berita.TanggalPelaksanaan,
			&berita.Deskripsi,
			&galeri.IdGaleri,
			&galeri.Gambar,
		)
		if err != nil {
			return models.Berita{}, err
		}

		// Hanya ambil data berita sekali (dari baris pertama)
		if isFirstRow {
			isFirstRow = false
		}

		// Tambahkan galeri jika ada
		if galeri.IdGaleri != "" {
			galeri.IdBerita = berita.IdBerita
			galeriList = append(galeriList, galeri)
		}
	}

	if err = rows.Err(); err != nil {
		return models.Berita{}, err
	}

	// Jika tidak ada baris yang ditemukan
	if berita.IdBerita == "" {
		return models.Berita{}, sql.ErrNoRows
	}

	// Set galeri ke berita
	berita.GambarBerita = galeriList

	return berita, nil
}

// UpdateBerita implements BeritaRepository.
func (b *beritaRepositoryImpl) UpdateBerita(ctx context.Context, tx *sql.Tx, berita models.Berita) error {
	panic("unimplemented")
}
