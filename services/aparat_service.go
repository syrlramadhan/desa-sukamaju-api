package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/models"
	"github.com/syrlramadhan/desa-sukamaju-api/repositories"
)

type AparatService interface {
	CreateAparat(ctx context.Context, r *http.Request, aparatReq dto.AparatRequest) (dto.AparatResponse, int, error)
	UpdateAparat(ctx context.Context, r *http.Request, idAparat string, aparatReq dto.AparatRequest) (dto.AparatResponse, int, error)
	GetAllAparat(ctx context.Context) ([]dto.AparatResponse, int, error)
	GetAparatById(ctx context.Context, idAparat string) (dto.AparatResponse, int, error)
	DeleteAparat(ctx context.Context, idAparat string) (int, error)
	BulkDeleteAparat(ctx context.Context, idAparat []string) (int, error)
}

type aparatServiceImpl struct {
	AparatRepository repositories.AparatRepository
	DB               *sql.DB
}

func NewAparatService(aparatRepository repositories.AparatRepository, db *sql.DB) AparatService {
	return &aparatServiceImpl{
		AparatRepository: aparatRepository,
		DB:               db,
	}
}

// CreateAparat implements AparatService.
func (a *aparatServiceImpl) CreateAparat(ctx context.Context, r *http.Request, aparatReq dto.AparatRequest) (dto.AparatResponse, int, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to parse form: %v", err)
	}

	aparatReq.Nama = r.FormValue("nama")
	aparatReq.Jabatan = r.FormValue("jabatan")
	aparatReq.NoTelepon = r.FormValue("no_telepon")
	aparatReq.Email = r.FormValue("email")
	aparatReq.Status = r.FormValue("status")
	aparatReq.PeriodeMulai = r.FormValue("periode_mulai")
	aparatReq.PeriodeSelesai = r.FormValue("periode_selesai")
	fotoFile, _, err := r.FormFile("foto")
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get foto from form: %v", err)
	}
	defer fotoFile.Close()

	uuid := uuid.New().String()

	fotoFilename := fmt.Sprintf("aparat_%s_%s.jpg", aparatReq.Nama, uuid)

	uploadDir := "uploads/aparat/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := filepath.Join(uploadDir, fotoFilename)
	out, err := os.Create(filePath)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, fotoFile)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to copy file: %v", err)
	}

	aparatReq.Foto = fotoFilename

	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	aparatModel := models.Aparat{
		IdAparat:       uuid,
		Nama:           aparatReq.Nama,
		Jabatan:        aparatReq.Jabatan,
		NoTelepon:      aparatReq.NoTelepon,
		Email:          aparatReq.Email,
		Status:         aparatReq.Status,
		PeriodeMulai:   aparatReq.PeriodeMulai,
		PeriodeSelesai: aparatReq.PeriodeSelesai,
		Foto:           aparatReq.Foto,
	}

	addAparat, err := a.AparatRepository.AddAparat(ctx, tx, aparatModel)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal menambahkan aparat: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return helpers.ConvertAparatToResponseDTO(addAparat), http.StatusOK, nil
}

// GetAllAparat implements AparatService.
func (a *aparatServiceImpl) GetAllAparat(ctx context.Context) ([]dto.AparatResponse, int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	aparats, err := a.AparatRepository.GetAllAparat(ctx, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan data aparat: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return helpers.ConvertAparatToListResDTO(aparats), http.StatusOK, nil
}

// GetAparatById implements AparatService.
func (a *aparatServiceImpl) GetAparatById(ctx context.Context, idAparat string) (dto.AparatResponse, int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	aparat, err := a.AparatRepository.GetAparatById(ctx, tx, idAparat)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan aparat: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return helpers.ConvertAparatToResponseDTO(aparat), http.StatusOK, nil
}

// UpdateAparat implements AparatService.
func (a *aparatServiceImpl) UpdateAparat(ctx context.Context, r *http.Request, idAparat string, aparatReq dto.AparatRequest) (dto.AparatResponse, int, error) {
	aparatReq.Nama = r.FormValue("nama")
	aparatReq.Jabatan = r.FormValue("jabatan")
	aparatReq.NoTelepon = r.FormValue("no_telepon")
	aparatReq.Email = r.FormValue("email")
	aparatReq.Status = r.FormValue("status")
	aparatReq.PeriodeMulai = r.FormValue("periode_mulai")
	aparatReq.PeriodeSelesai = r.FormValue("periode_selesai")
	fotoFile, _, err := r.FormFile("foto")
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get foto from form: %v", err)
	}
	defer fotoFile.Close()

	fotoFilename := fmt.Sprintf("aparat_%s_%s.jpg", aparatReq.Nama, idAparat)

	uploadDir := "uploads/aparat/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := filepath.Join(uploadDir, fotoFilename)
	out, err := os.Create(filePath)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, fotoFile)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to copy file: %v", err)
	}

	aparatReq.Foto = fotoFilename

	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	getAparat, err := a.AparatRepository.GetAparatById(ctx, tx, idAparat)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan aparat: %v", err)
	}

	if getAparat.IdAparat == "" {
		return dto.AparatResponse{}, http.StatusNotFound, fmt.Errorf("aparat dengan id %s tidak ditemukan", idAparat)
	}

	if _, err := os.Stat(uploadDir); os.IsExist(err) {
		os.Remove(filepath.Join(uploadDir, getAparat.Foto))
	}

	aparatModel := models.Aparat{
		IdAparat:       idAparat,
		Nama:           aparatReq.Nama,
		Jabatan:        aparatReq.Jabatan,
		NoTelepon:      aparatReq.NoTelepon,
		Email:          aparatReq.Email,
		Status:         aparatReq.Status,
		PeriodeMulai:   aparatReq.PeriodeMulai,
		PeriodeSelesai: aparatReq.PeriodeSelesai,
		Foto:           aparatReq.Foto,
	}

	err = a.AparatRepository.UpdateAparat(ctx, tx, aparatModel)
	if err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengupdate aparat: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return dto.AparatResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return helpers.ConvertAparatToResponseDTO(aparatModel), http.StatusOK, nil
}

// DeleteAparat implements AparatService.
func (a *aparatServiceImpl) DeleteAparat(ctx context.Context, idAparat string) (int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	getAparat, err := a.AparatRepository.GetAparatById(ctx, tx, idAparat)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan aparat: %v", err)
	}

	if getAparat.IdAparat == "" {
		return http.StatusNotFound, fmt.Errorf("aparat dengan id %s tidak ditemukan", idAparat)
	}

	err = a.AparatRepository.DeleteAparat(ctx, tx, idAparat)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus aparat: %v", err)
	}

	uploadDir := "uploads/aparat/"
	if _, err := os.Stat(uploadDir); os.IsExist(err) {
		os.Remove(filepath.Join(uploadDir, getAparat.Foto))
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}

// BulkDeleteAparat implements AparatService.
func (a *aparatServiceImpl) BulkDeleteAparat(ctx context.Context, idAparat []string) (int, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	var aparat []models.Aparat

	for _, id := range idAparat {
		getAparat, err := a.AparatRepository.GetAparatById(ctx, tx, id)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan aparat: %v", err)
		}

		if getAparat.IdAparat == "" {
			return http.StatusNotFound, fmt.Errorf("aparat dengan id %s tidak ditemukan", id)
		}

		aparat = append(aparat, getAparat)
	}

	err = a.AparatRepository.BulkDeleteAparat(ctx, tx, idAparat)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus aparat: %v", err)
	}

	for _, aparat := range aparat {
		uploadDir := "uploads/aparat/"
		if _, err := os.Stat(uploadDir); os.IsExist(err) {
			os.Remove(filepath.Join(uploadDir, aparat.Foto))
		}
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}
