package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

type AdminRepository interface {
	GetAdmin(ctx context.Context, tx *sql.Tx, username string) (models.Admin, error)
	GetAdminById(ctx context.Context, tx *sql.Tx, idAdmin string) (models.Admin, error)
	UpdateUsernameAdmin(ctx context.Context, tx *sql.Tx, oldUsername, newUsername string) error
	UpdatePasswordAdmin(ctx context.Context, tx *sql.Tx, idAdmin, newPassword string) error
}

type AdminRepositoryImpl struct {
}

func NewAdminRepository() AdminRepository {
	return &AdminRepositoryImpl{}
}

// GetAdmin implements AdminRepository.
func (a *AdminRepositoryImpl) GetAdmin(ctx context.Context, tx *sql.Tx, username string) (models.Admin, error) {
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

// GetAdminById implements AdminRepository.
func (a *AdminRepositoryImpl) GetAdminById(ctx context.Context, tx *sql.Tx, idAdmin string) (models.Admin, error) {
	var admin models.Admin
	query := "SELECT id_admin, username, password FROM admin WHERE id_admin = ?"

	err := tx.QueryRowContext(ctx, query, idAdmin).Scan(
		&admin.IdAdmin,
		&admin.Username,
		&admin.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Admin{}, errors.New("id tidak ditemukan")
		}
		return models.Admin{}, err
	}

	return admin, nil
}

// UpdateUsernameAdmin implements AdminRepository.
func (a *AdminRepositoryImpl) UpdateUsernameAdmin(ctx context.Context, tx *sql.Tx, oldUsername string, newUsername string) error {
	// Update username
	updateQuery := "UPDATE admin SET username = ? WHERE username = ?"
	result, err := tx.ExecContext(ctx, updateQuery, newUsername, oldUsername)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("admin tidak ditemukan")
	}

	return nil
}

// UpdatePasswordAdmin implements AdminRepository.
func (a *AdminRepositoryImpl) UpdatePasswordAdmin(ctx context.Context, tx *sql.Tx, idAdmin string, newPassword string) error {
	updateQuery := "UPDATE admin SET password = ? WHERE id_admin = ?"
	result, err := tx.ExecContext(ctx, updateQuery, newPassword, idAdmin)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("admin tidak ditemukan")
	}

	return nil
}
