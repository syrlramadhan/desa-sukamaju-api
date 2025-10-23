package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/services"
)

type PendudukController interface {
	GetPenduduk(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UpdatePenduduk(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type pendudukControllerImpl struct {
	service services.PendudukService
}

func NewPendudukController(service services.PendudukService) PendudukController {
	return &pendudukControllerImpl{
		service: service,
	}
}

// GetPenduduk implements PendudukController.
func (p *pendudukControllerImpl) GetPenduduk(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pendudukResponse, statusCode, err := p.service.GetPenduduk(r.Context())
	if err != nil {
		helpers.WriteJSONError(w, statusCode, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, pendudukResponse, "berhasil mendapatkan data penduduk")
}

// UpdatePenduduk implements PendudukController.
func (p *pendudukControllerImpl) UpdatePenduduk(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pendudukReq := dto.PendudukRequest{}
	helpers.ReadFromRequestBody(r, &pendudukReq)
	
	pendudukResponse, statusCode, err := p.service.UpdatePenduduk(r.Context(), pendudukReq)
	if err != nil {
		helpers.WriteJSONError(w, statusCode, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, pendudukResponse, "berhasil memperbarui data penduduk")
}