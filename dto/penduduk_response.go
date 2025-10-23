package dto

type PendudukResponse struct {
	IdPenduduk          string `json:"id_penduduk"`
	TotalPenduduk       string `json:"total_penduduk"`
	TotalKepalaKeluarga string `json:"total_kepala_keluarga"`
}