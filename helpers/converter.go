package helpers

import (
	"github.com/syrlramadhan/desa-sukamaju-api/dto"
	"github.com/syrlramadhan/desa-sukamaju-api/models"
)

func ConvertAparatToResponseDTO(aparat models.Aparat) dto.AparatResponse {
	return dto.AparatResponse{
		IdAparat:       aparat.IdAparat,
		Nama:           aparat.Nama,
		Jabatan:        aparat.Jabatan,
		NoTelepon:      aparat.NoTelepon,
		Email:          aparat.Email,
		Status:         aparat.Status,
		PeriodeMulai:   aparat.PeriodeMulai,
		PeriodeSelesai: aparat.PeriodeSelesai,
		Foto:           aparat.Foto,
	}
}

func ConvertAparatToListResDTO(aparats []models.Aparat) []dto.AparatResponse {
	var aparatResponse []dto.AparatResponse

	for _, aparat := range aparats {
		aparatResponse = append(aparatResponse, ConvertAparatToResponseDTO(aparat))
	}

	return aparatResponse
}