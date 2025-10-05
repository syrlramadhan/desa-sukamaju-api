package repositories

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

type KontenRepository interface {
	GetKontak(ctx context.Context, tx *sql.Tx, kontak models.Kontak) (models.Kontak, error)
	UpdateKontak(ctx context.Context, tx *sql.Tx, kontak models.Kontak) error
}

type kontenRepositoryImpl struct {
}

func NewKontenRepository() KontenRepository {
	return &kontenRepositoryImpl{}
}

// GetKontak implements KontenRepository.
func (k *kontenRepositoryImpl) GetKontak(ctx context.Context, tx *sql.Tx, kontak models.Kontak) (models.Kontak, error) {
	query := "SELECT id_kontak, email, telepon, facebook, instagram, youtube FROM kontak WHERE id_kontak = ?"

	tx.QueryRowContext(ctx, query, kontak.IdKontak).Scan(
		&kontak.IdKontak,
		&kontak.Email,
		&kontak.Telepon,
		&kontak.Facebook,
		&kontak.Instagram,
		&kontak.Youtube,
	)

	return kontak, nil
}

// UpdateKontak implements KontenRepository.
func (k *kontenRepositoryImpl) UpdateKontak(ctx context.Context, tx *sql.Tx, kontak models.Kontak) error {
	query := "UPDATE kontak SET email = ?, telepon = ?, facebook = ?, instagram = ?, youtube = ? WHERE id_kontak = ?"

	_, err := tx.ExecContext(ctx, query, kontak.Email, kontak.Telepon, kontak.Facebook, kontak.Instagram, kontak.Youtube, kontak.IdKontak)
	
	return err
}
