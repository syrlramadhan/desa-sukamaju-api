package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/repositories"
)

type AdminService interface {
	GetAdminById(ctx context.Context, idAdmin string) (dto.AdminResponse, int, error)
	LoginAdmin(ctx context.Context, loginReq dto.LoginRequest) (string, int, error)
	UpdateUsernameAdmin(ctx context.Context, idAdmin string, usernameReq dto.UpdateUsernameRequest) (int, error)
	UpdatePasswordAdmin(ctx context.Context, username string, passReq dto.UpdatePasswordRequest) (int, error)
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

// GetAdminById implements AdminService.
func (a *AdminServiceImpl) GetAdminById(ctx context.Context, idAdmin string) (dto.AdminResponse, int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AdminResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	admin, err := a.AdminRepository.GetAdminById(ctx, tx, idAdmin)
	if err != nil {
		return dto.AdminResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan admin: %v", err)
	}

	adminResponse := dto.AdminResponse{
		IdAdmin:  admin.IdAdmin,
		Username: admin.Username,
	}

	if err := tx.Commit(); err != nil {
		return dto.AdminResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return adminResponse, http.StatusOK, nil
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

	admin, err := a.AdminRepository.GetAdmin(ctx, tx, loginReq.Username)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan admin: %v", err)
	}

	if !helpers.VerifyPassword(admin.Password, loginReq.Password) {
		return "", http.StatusBadRequest, fmt.Errorf("username atau password salah")
	}

	token, err := helpers.GenerateJWT(admin.IdAdmin, admin.Username)
	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("gagal menghasilkan token: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return token, http.StatusOK, nil
}

// UpdatePasswordAdmin implements AdminService.
func (a *AdminServiceImpl) UpdatePasswordAdmin(ctx context.Context, username string, passReq dto.UpdatePasswordRequest) (int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	admin, err := a.AdminRepository.GetAdmin(ctx, tx, username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan admin: %v", err)
	}

	if passReq.OldPassword == "" {
		return http.StatusBadRequest, fmt.Errorf("password lama tidak boleh kosong")
	}

	if !helpers.VerifyPassword(admin.Password, passReq.OldPassword) {
		return http.StatusBadRequest, fmt.Errorf("password lama salah")
	}

	if passReq.Password == "" {
		return http.StatusBadRequest, fmt.Errorf("password tidak boleh kosong")
	}

	if passReq.KonfirmasiPassword == "" {
		return http.StatusBadRequest, fmt.Errorf("konfirmasi password tidak boleh kosong")
	}

	if passReq.Password != passReq.KonfirmasiPassword {
		return http.StatusBadRequest, fmt.Errorf("password dan konfirmasi password tidak cocok")
	}

	hashedPassword, err := helpers.HashPassword(passReq.Password)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal hash password: %v", err)
	}

	err = a.AdminRepository.UpdatePasswordAdmin(ctx, tx, admin.IdAdmin, hashedPassword)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui password: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}

// UpdateUsernameAdmin implements AdminService.
func (a *AdminServiceImpl) UpdateUsernameAdmin(ctx context.Context, oldUsername string, usernameReq dto.UpdateUsernameRequest) (int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	admin, err := a.AdminRepository.GetAdmin(ctx, tx, oldUsername)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan admin: %v", err)
	}

	if usernameReq.Username == "" {
		return http.StatusBadRequest, fmt.Errorf("username tidak boleh kosong")
	}

	if admin.Username == usernameReq.Username {
		return http.StatusBadRequest, fmt.Errorf("username tidak boleh sama dengan username lama")
	}

	err = a.AdminRepository.UpdateUsernameAdmin(ctx, tx, oldUsername, usernameReq.Username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui username: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}
