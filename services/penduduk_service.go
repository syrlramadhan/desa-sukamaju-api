package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/models"
	"github.com/syrlramadhan/desa-sukamaju-api/repositories"
)

type PendudukService interface {
	GetPenduduk(ctx context.Context) (dto.PendudukResponse, int, error)
	UpdatePenduduk(ctx context.Context, pendudukRequest dto.PendudukRequest) (dto.PendudukResponse, int, error)
}

type pendudukServiceImpl struct {
	repo repositories.PendudukRepository
	DB   *sql.DB
}

func NewPendudukService(repo repositories.PendudukRepository, DB *sql.DB) PendudukService {
	return &pendudukServiceImpl{
		repo: repo,
		DB:   DB,
	}
}

// GetPenduduk implements PendudukService.
func (p *pendudukServiceImpl) GetPenduduk(ctx context.Context) (dto.PendudukResponse, int, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return dto.PendudukResponse{}, 500, err
	}
	defer tx.Rollback()

	pendudukModel := models.Penduduk{}
	penduduk, err := p.repo.GetPenduduk(ctx, tx, pendudukModel)
	if err != nil {
		return dto.PendudukResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan data penduduk: %v", err)
	}

	pendudukResponse := dto.PendudukResponse{
		IdPenduduk:          penduduk.IdPenduduk,
		TotalPenduduk:       penduduk.TotalPenduduk,
		TotalKepalaKeluarga: penduduk.TotalKepalaKeluarga,
	}

	if err := tx.Commit(); err != nil {
		return dto.PendudukResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return pendudukResponse, http.StatusOK, nil
}

// UpdatePenduduk implements PendudukService.
func (p *pendudukServiceImpl) UpdatePenduduk(ctx context.Context, pendudukRequest dto.PendudukRequest) (dto.PendudukResponse, int, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return dto.PendudukResponse{}, 500, err
	}
	defer tx.Rollback()

	pendudukModel := models.Penduduk{
		TotalPenduduk:       pendudukRequest.TotalPenduduk,
		TotalKepalaKeluarga: pendudukRequest.TotalKepalaKeluarga,
	}

	penduduk, err := p.repo.GetPenduduk(ctx, tx, pendudukModel)
	if err != nil {
		return dto.PendudukResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan data penduduk: %v", err)
	}

	pendudukModel.IdPenduduk = penduduk.IdPenduduk

	err = p.repo.UpdatePenduduk(ctx, tx, pendudukModel)
	if err != nil {
		return dto.PendudukResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memperbarui data penduduk: %v", err)
	}

	updatedPendudukResponse := dto.PendudukResponse{
		IdPenduduk:          pendudukModel.IdPenduduk,
		TotalPenduduk:       pendudukModel.TotalPenduduk,
		TotalKepalaKeluarga: pendudukModel.TotalKepalaKeluarga,
	}

	if err := tx.Commit(); err != nil {
		return dto.PendudukResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return updatedPendudukResponse, http.StatusOK, nil
}
