package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/helpers"
	"github.com/syrlramadhan/desa-sukamaju-api/services"
)

type AdminController interface {
	GetAdminById(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	LoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UpdatePasswordAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateUsernameAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type AdminControllerImpl struct {
	AdminService services.AdminService
}

func NewAdminController(adminService services.AdminService) AdminController {
	return &AdminControllerImpl{
		AdminService: adminService,
	}
}

// GetAdminById implements AdminController.
func (a *AdminControllerImpl) GetAdminById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	adminResponse, code, err := a.AdminService.GetAdminById(r.Context(), ps.ByName("id_admin"))
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, adminResponse, "berhasil mendapatkan data admin")
}

// LoginAdmin implements AdminController.
func (a *AdminControllerImpl) LoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginReq := dto.LoginRequest{}
	helpers.ReadFromRequestBody(r, &loginReq)

	token, code, err := a.AdminService.LoginAdmin(r.Context(), loginReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONLogin(w, token, "berhasil login")
}

// UpdatePasswordAdmin implements AdminController.
func (a *AdminControllerImpl) UpdatePasswordAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	passReq := dto.UpdatePasswordRequest{}
	helpers.ReadFromRequestBody(r, &passReq)

	code, err := a.AdminService.UpdatePasswordAdmin(r.Context(), ps.ByName("username"), passReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "password berhasil diperbarui")
}

// UpdateUsernameAdmin implements AdminController.
func (a *AdminControllerImpl) UpdateUsernameAdmin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	usernameReq := dto.UpdateUsernameRequest{}
	helpers.ReadFromRequestBody(r, &usernameReq)

	code, err := a.AdminService.UpdateUsernameAdmin(r.Context(), ps.ByName("username"), usernameReq)
	if err != nil {
		helpers.WriteJSONError(w, code, err.Error())
		return
	}

	helpers.WriteJSONNoData(w, "username berhasil diperbarui")
}
