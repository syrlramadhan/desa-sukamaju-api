package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/services"
)

type KontenController interface {
	GetKontak(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateKontak(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type kontenControllerImpl struct {
	KontenService services.KontenService
}

func NewKontenController(kontenService services.KontenService) KontenController {
	return &kontenControllerImpl{
		KontenService: kontenService,
	}
}

// GetKontak implements KontenController.
func (k *kontenControllerImpl) GetKontak(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	responseDTO, code, err := k.KontenService.GetKontak(r.Context(), ps.ByName("id_kontak"))
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan data kontak")
}

// UpdateKontak implements KontenController.
func (k *kontenControllerImpl) UpdateKontak(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	kontakReq := dto.KontakRequest{}
	helpers.ReadFromRequestBody(r, &kontakReq)

	code, err := k.KontenService.UpdateKontak(r.Context(), ps.ByName("id_kontak"), kontakReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "berhasil memperbarui data kontak")
}
