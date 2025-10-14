package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

type AparatRepository interface {
	AddAparat(ctx context.Context, tx *sql.Tx, aparat models.Aparat) (models.Aparat, error)
	GetAllAparat(ctx context.Context, tx *sql.Tx) ([]models.Aparat, error)
	GetAparatById(ctx context.Context, tx *sql.Tx, idAparat string) (models.Aparat, error)
	UpdateAparat(ctx context.Context, tx *sql.Tx, aparat models.Aparat) error
	DeleteAparat(ctx context.Context, tx *sql.Tx, idAparat string) error
	BulkDeleteAparat(ctx context.Context, tx *sql.Tx, idAparat []string) error
}

type aparatRepositoryImpl struct {
}

func NewAparatRepository() AparatRepository {
	return &aparatRepositoryImpl{}
}

// AddAparat implements AparatRepository.
func (a *aparatRepositoryImpl) AddAparat(ctx context.Context, tx *sql.Tx, aparat models.Aparat) (models.Aparat, error) {
	query := "INSERT INTO aparat (id_aparat, nama, jabatan, no_telepon, email, status, periode_mulai, periode_selesai, foto) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, aparat.IdAparat, aparat.Nama, aparat.Jabatan, aparat.NoTelepon, aparat.Email, aparat.Status, aparat.PeriodeMulai, aparat.PeriodeSelesai, aparat.Foto)
	if err != nil {
		return models.Aparat{}, err
	}

	return aparat, nil
}

// GetAllAparat implements AparatRepository.
func (a *aparatRepositoryImpl) GetAllAparat(ctx context.Context, tx *sql.Tx) ([]models.Aparat, error) {
	query := "SELECT id_aparat, nama, jabatan, no_telepon, email, status, periode_mulai, periode_selesai, foto FROM aparat"

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aparats []models.Aparat
	for rows.Next() {
		var aparat models.Aparat
		if err := rows.Scan(&aparat.IdAparat, &aparat.Nama, &aparat.Jabatan, &aparat.NoTelepon, &aparat.Email, &aparat.Status, &aparat.PeriodeMulai, &aparat.PeriodeSelesai, &aparat.Foto); err != nil {
			return nil, err
		}
		aparats = append(aparats, aparat)
	}

	return aparats, nil
}

// GetAparatById implements AparatRepository.
func (a *aparatRepositoryImpl) GetAparatById(ctx context.Context, tx *sql.Tx, idAparat string) (models.Aparat, error) {
	query := "SELECT id_aparat, nama, jabatan, no_telepon, email, status, periode_mulai, periode_selesai, foto FROM aparat WHERE id_aparat = ?"

	var aparat models.Aparat
	err := tx.QueryRowContext(ctx, query, idAparat).Scan(&aparat.IdAparat, &aparat.Nama, &aparat.Jabatan, &aparat.NoTelepon, &aparat.Email, &aparat.Status, &aparat.PeriodeMulai, &aparat.PeriodeSelesai, &aparat.Foto)
	if err != nil {
		return models.Aparat{}, err
	}

	return aparat, nil
}

// UpdateAparat implements AparatRepository.
func (a *aparatRepositoryImpl) UpdateAparat(ctx context.Context, tx *sql.Tx, aparat models.Aparat) error {
	query := "UPDATE aparat SET nama = ?, jabatan = ?, no_telepon = ?, email = ?, status = ?, periode_mulai = ?, periode_selesai = ?, foto = ? WHERE id_aparat = ?"

	_, err := tx.ExecContext(ctx, query, aparat.Nama, aparat.Jabatan, aparat.NoTelepon, aparat.Email, aparat.Status, aparat.PeriodeMulai, aparat.PeriodeSelesai, aparat.Foto, aparat.IdAparat)
	return err
}

// DeleteAparat implements AparatRepository.
func (a *aparatRepositoryImpl) DeleteAparat(ctx context.Context, tx *sql.Tx, idAparat string) error {
	query := "DELETE FROM aparat WHERE id_aparat = ?"

	_, err := tx.ExecContext(ctx, query, idAparat)
	return err
}

// BulkDeleteAparat implements AparatRepository.
func (a *aparatRepositoryImpl) BulkDeleteAparat(ctx context.Context, tx *sql.Tx, idAparat []string) error {
	if len(idAparat) == 0 {
		return nil // Tidak ada yang perlu dihapus
	}

	// Buat placeholder untuk IN clause (?, ?, ?, ...)
	placeholders := make([]string, len(idAparat))
	args := make([]interface{}, len(idAparat))

	for i, id := range idAparat {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM aparat WHERE id_aparat IN (%s)",
		strings.Join(placeholders, ","))

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
