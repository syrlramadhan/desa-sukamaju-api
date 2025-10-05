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

type KontenService interface {
	GetKontak(ctx context.Context, idKontak string) (dto.KontakResponse, int, error)
	UpdateKontak(ctx context.Context, idKontak string, kontakReq dto.KontakRequest) (int, error)
}

type kontenServiceImpl struct {
	KontenRepository repositories.KontenRepository
	DB               *sql.DB
}

func NewKontenService(kontenRepository repositories.KontenRepository, db *sql.DB) KontenService {
	return &kontenServiceImpl{
		KontenRepository: kontenRepository,
		DB:               db,
	}
}

// GetKontak implements KontenService.
func (k *kontenServiceImpl) GetKontak(ctx context.Context, idKontak string) (dto.KontakResponse, int, error) {
	tx, err := k.DB.Begin()
	if err != nil {
		return dto.KontakResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	kontakModel := models.Kontak{
		IdKontak: idKontak,
	}

	kontak, err := k.KontenRepository.GetKontak(ctx, tx, kontakModel)
	if err != nil {
		return dto.KontakResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mendapatkan kontak: %v", err)
	}

	kontakResponse := dto.KontakResponse{
		IdKontak:  kontak.IdKontak,
		Email:     kontak.Email,
		Telepon:   kontak.Telepon,
		Facebook:  kontak.Facebook,
		Youtube:   kontak.Youtube,
		Instagram: kontak.Instagram,
	}

	if err := tx.Commit(); err != nil {
		return dto.KontakResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return kontakResponse, http.StatusOK, nil
}

// UpdateKontak implements KontenService.
func (k *kontenServiceImpl) UpdateKontak(ctx context.Context, idKontak string, kontakReq dto.KontakRequest) (int, error) {
	tx, err := k.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memulai transaksi: %v", err)
	}

	kontakModel := models.Kontak{
		IdKontak:  idKontak,
		Email:     kontakReq.Email,
		Telepon:   kontakReq.Telepon,
		Facebook:  kontakReq.Facebook,
		Instagram: kontakReq.Instagram,
		Youtube:   kontakReq.Youtube,
	}

	err = k.KontenRepository.UpdateKontak(ctx, tx, kontakModel)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengupdate kontak: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengkomit transaksi: %v", err)
	}

	return http.StatusOK, nil
}
