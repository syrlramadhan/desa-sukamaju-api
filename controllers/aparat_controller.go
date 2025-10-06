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
	aparatResponse , code, err := a.AparatService.GetAllAparat(r.Context())
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, aparatResponse, "berhasil mendapatkan data aparat")
}

// GetAparatById implements AparatController.
func (a *aparatControllerImpl) GetAparatById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	panic("unimplemented")
}

// UpdateAparat implements AparatController.
func (a *aparatControllerImpl) UpdateAparat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	panic("unimplemented")
}

// DeleteAparat implements AparatController.
func (a *aparatControllerImpl) DeleteAparat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	panic("unimplemented")
}