package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	helper "github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/services"
)

type AdminController interface {
	LoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UpdateAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type AdminControllerImpl struct {
	AdminService services.AdminService
}

func NewAdminController(adminService services.AdminService) AdminController {
	return &AdminControllerImpl{
		AdminService: adminService,
	}
}

// LoginAdmin implements AdminController.
func (a *AdminControllerImpl) LoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginReq := dto.LoginRequest{}
	helper.ReadFromRequestBody(r, &loginReq)

	token, code, err := a.AdminService.LoginAdmin(r.Context(), loginReq)
	if err != nil {
		helper.WriteJSONError(w, code, err.Error())
		return
	}

	helper.WriteJSONLogin(w, token, "berhasil login")
}

// UpdateAdmin implements AdminController.
func (a *AdminControllerImpl) UpdateAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	panic("unimplemented")
}
