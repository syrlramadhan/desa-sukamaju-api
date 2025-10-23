package repositories

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

type PendudukRepository interface {
	GetPenduduk(ctx context.Context, tx *sql.Tx, penduduk models.Penduduk) (models.Penduduk, error)
	UpdatePenduduk(ctx context.Context, tx *sql.Tx, penduduk models.Penduduk) error
}

type pendudukRepositoryImpl struct {
}

func NewPendudukRepository() PendudukRepository {
	return &pendudukRepositoryImpl{}
}

// GetPenduduk implements PendudukRepository.
func (p *pendudukRepositoryImpl) GetPenduduk(ctx context.Context, tx *sql.Tx, penduduk models.Penduduk) (models.Penduduk, error) {
	query := "SELECT id_penduduk, total_penduduk, total_kepala_keluarga FROM penduduk"
	row := tx.QueryRowContext(ctx, query)
	if err := row.Scan(&penduduk.IdPenduduk, &penduduk.TotalPenduduk, &penduduk.TotalKepalaKeluarga); err != nil {
		return models.Penduduk{}, err
	}
	return penduduk, nil
}

// UpdatePenduduk implements PendudukRepository.
func (p *pendudukRepositoryImpl) UpdatePenduduk(ctx context.Context, tx *sql.Tx, penduduk models.Penduduk) error {
	query := "UPDATE penduduk SET total_penduduk = ?, total_kepala_keluarga = ? WHERE id_penduduk = ?"
	_, err := tx.ExecContext(ctx, query, penduduk.TotalPenduduk, penduduk.TotalKepalaKeluarga, penduduk.IdPenduduk)
	return err
}