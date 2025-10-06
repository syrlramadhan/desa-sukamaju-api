package repositories

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

type AparatRepository interface {
	AddAparat(ctx context.Context, tx *sql.Tx, aparat models.Aparat) (models.Aparat, error)
	GetAllAparat(ctx context.Context, tx *sql.Tx) ([]models.Aparat, error)
	GetAparatById(ctx context.Context, tx *sql.Tx, idAparat string) (models.Aparat, error)
	UpdateAparat(ctx context.Context, tx *sql.Tx, aparat models.Aparat) error
	DeleteAparat(ctx context.Context, tx *sql.Tx, idAparat string) error
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
	panic("unimplemented")
}

// UpdateAparat implements AparatRepository.
func (a *aparatRepositoryImpl) UpdateAparat(ctx context.Context, tx *sql.Tx, aparat models.Aparat) error {
	panic("unimplemented")
}

// DeleteAparat implements AparatRepository.
func (a *aparatRepositoryImpl) DeleteAparat(ctx context.Context, tx *sql.Tx, idAparat string) error {
	panic("unimplemented")
}