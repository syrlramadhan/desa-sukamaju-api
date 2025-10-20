package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/services"
)

type BeritaController interface {
	CreateBerita(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	CreatePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetAllBerita(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetBeritaById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateBerita(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	DeleteBerita(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type beritaControllerImpl struct {
	BeritaService services.BeritaService
}

func NewBeritaController(beritaService services.BeritaService) BeritaController {
	return &beritaControllerImpl{
		BeritaService: beritaService,
	}
}

// CreateBerita implements BeritaController.
func (b *beritaControllerImpl) CreateBerita(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beritaReq := dto.BeritaRequest{}

	code, err := b.BeritaService.CreateBerita(r.Context(), r, beritaReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil menambahkan berita")
}

// CreatePhoto implements BeritaController.
func (b *beritaControllerImpl) CreatePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	photoReq := dto.GaleriRequest{}

	code, err := b.BeritaService.CreatePhoto(r.Context(), r, photoReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil menambahkan photo berita")
}

// DeleteBerita implements BeritaController.
func (b *beritaControllerImpl) DeleteBerita(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	panic("unimplemented")
}

// GetAllBerita implements BeritaController.
func (b *beritaControllerImpl) GetAllBerita(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseDTO, code, err := b.BeritaService.GetAllBerita(r.Context())
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan data berita")
}

// GetBeritaById implements BeritaController.
func (b *beritaControllerImpl) GetBeritaById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	responseDTO, code, err := b.BeritaService.GetBeritaById(r.Context(), ps.ByName("id_berita"))
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan data berita")
}

// UpdateBerita implements BeritaController.
func (b *beritaControllerImpl) UpdateBerita(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("unimplemented")
}
