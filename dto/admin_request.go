package dto

type UpdateUsernameRequest struct {
	Username string `json:"username"`
}

type UpdatePasswordRequest struct {
	OldPassword        string `json:"old_password"`
	Password           string `json:"password"`
	KonfirmasiPassword string `json:"konfirmasi_password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
