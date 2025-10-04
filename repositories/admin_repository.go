package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

type AdminRepository interface {
	GetAdmin(ctx context.Context, tx *sql.Tx, username, password string) (models.Admin,error)
	UpdateAdmin(ctx context.Context, tx *sql.Tx, idAdmin, username, password string) error
}

type AdminRepositoryImpl struct {
}

func NewAdminRepository() AdminRepository {
	return &AdminRepositoryImpl{}
}

// GetAdmin implements AdminRepository.
func (a *AdminRepositoryImpl) GetAdmin(ctx context.Context, tx *sql.Tx, username string, password string) (models.Admin,error) {
	var admin models.Admin
	query := "SELECT id_admin, username, password FROM admin WHERE username = ?"

	err := tx.QueryRowContext(ctx, query, username).Scan(
		&admin.IdAdmin,
		&admin.Username,
		&admin.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Admin{}, errors.New("username tidak ditemukan")
		}
		return models.Admin{}, err
	}

	return admin, nil
}

// UpdateAdmin implements AdminRepository.
func (a *AdminRepositoryImpl) UpdateAdmin(ctx context.Context, tx *sql.Tx, idAdmin string, username string, password string) error {
	panic("unimplemented")
}