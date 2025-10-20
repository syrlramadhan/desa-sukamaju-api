package controllers

import (
	"encoding/json"
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
	DeletePhotoByFilename(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	BulkDeletePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
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
	idBerita := ps.ByName("id_berita")

	code, err := b.BeritaService.DeleteBerita(r.Context(), idBerita)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil menghapus berita")
}

// DeletePhotoByFilename implements BeritaController.
func (b *beritaControllerImpl) DeletePhotoByFilename(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code, err := b.BeritaService.DeletePhotoByFilename(r.Context(), ps.ByName("filename"))
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil menghapus foto")
}

// BulkDeletePhoto implements BeritaController.
func (b *beritaControllerImpl) BulkDeletePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request struct {
		Filenames []string `json:"filenames"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	code, err := b.BeritaService.BulkDeletePhoto(r.Context(), request.Filenames)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil menghapus foto")
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
func (b *beritaControllerImpl) UpdateBerita(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beritaReq := dto.BeritaRequest{}
	helpers.ReadFromRequestBody(r, &beritaReq)

	code, err := b.BeritaService.UpdateBerita(r.Context(), ps.ByName("id_berita"), beritaReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil mengupdate berita")
}