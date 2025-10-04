package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	helper "github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/repositories"
)

type AdminService interface {
	LoginAdmin(ctx context.Context, loginReq dto.LoginRequest) (string, int, error)
	UpdateAdmin(ctx context.Context, idAdmin, username, password string) (int, error)
}

type AdminServiceImpl struct {
	AdminRepository repositories.AdminRepository
	DB              *sql.DB
}

func NewAdminService(adminRepository repositories.AdminRepository, db *sql.DB) AdminService {
	return &AdminServiceImpl{
		AdminRepository: adminRepository,
		DB:              db,
	}
}

// LoginAdmin implements AdminService.
func (a *AdminServiceImpl) LoginAdmin(ctx context.Context, loginReq dto.LoginRequest) (string, int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	if loginReq.Username == "" {
		return "", http.StatusBadRequest, fmt.Errorf("username tidak boleh kosong")
	} else if loginReq.Password == "" {
		return "", http.StatusBadRequest, fmt.Errorf("password tidak boleh kosong")
	}

	admin, err := a.AdminRepository.GetAdmin(ctx, tx, loginReq.Username, loginReq.Password)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan admin: %v", err)
	}

	if !helper.VerifyPassword(admin.Password, loginReq.Password) {
		return "", http.StatusBadRequest, fmt.Errorf("username atau password salah")
	}

	token, err := helper.GenerateJWT(admin.Username)
	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("gagal menghasilkan token: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return token, http.StatusOK, nil
}

// UpdateAdmin implements AdminService.
func (a *AdminServiceImpl) UpdateAdmin(ctx context.Context, idAdmin string, username string, password string) (int, error) {
	panic("unimplemented")
}
