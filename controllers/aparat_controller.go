package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/services"
)

type AparatController interface {
	CreateAparat(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetAllAparat(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetAparatById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateAparat(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteAparat(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	BulkDeleteAparat(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type aparatControllerImpl struct {
	AparatService services.AparatService
}

func NewAparatController(aparatService services.AparatService) AparatController {
	return &aparatControllerImpl{
		AparatService: aparatService,
	}
}

// CreateAparat implements AparatController.
func (a *aparatControllerImpl) CreateAparat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	aparatReq := dto.AparatRequest{}

	aparatResponse, code, err := a.AparatService.CreateAparat(r.Context(), r, aparatReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, aparatResponse, "berhasil menambahkan data aparat")
}

// GetAllAparat implements AparatController.
func (a *aparatControllerImpl) GetAllAparat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	aparatResponse, code, err := a.AparatService.GetAllAparat(r.Context())
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, aparatResponse, "berhasil mendapatkan data aparat")
}

// GetAparatById implements AparatController.
func (a *aparatControllerImpl) GetAparatById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	aparatResponse, code, err := a.AparatService.GetAparatById(r.Context(), ps.ByName("id_aparat"))
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, aparatResponse, "berhasil mendapatkan data aparat")
}

// UpdateAparat implements AparatController.
func (a *aparatControllerImpl) UpdateAparat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	aparatReq := dto.AparatRequest{}

	aparatResponse, code, err := a.AparatService.UpdateAparat(r.Context(), r, ps.ByName("id_aparat"), aparatReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, aparatResponse, "berhasil memperbarui data aparat")
}

// DeleteAparat implements AparatController.
func (a *aparatControllerImpl) DeleteAparat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code, err := a.AparatService.DeleteAparat(r.Context(), ps.ByName("id_aparat"))
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil menghapus data aparat")
}

// BulkDeleteAparat implements AparatController.
func (a *aparatControllerImpl) BulkDeleteAparat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	idAparat := dto.BulkDeleteAparatRequest{}
	helpers.ReadFromRequestBody(r, &idAparat)

	code, err := a.AparatService.BulkDeleteAparat(r.Context(), idAparat.IDAparat)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil menghapus data aparat")
}
