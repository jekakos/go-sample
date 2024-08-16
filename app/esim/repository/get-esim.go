package repository

import (
	"sample-service/app/common"
	"sample-service/app/esim/dto"
	"sample-service/app/infra"
	"sample-service/app/model"
)

func GetEsim(params dto.GetEsimRequestDTO) ([]dto.GetEsimDbDTO, error) {

	var esims []model.Esim
	query := infra.DB.Select("id, iccid, lpa, user_uuid, profile_staus")

	iccid := common.GetStr(params.Iccid)
	userUuid := common.GetStr(params.UserUUID)

	if iccid != "" && userUuid != "" {
		query = query.Where("iccid = ? AND user_uuid = ?", common.ClearQuery(iccid), common.ClearQuery(userUuid))
	} else if iccid != "" {
		query = query.Where("iccid = ?", common.ClearQuery(iccid))
	} else if userUuid != "" {
		query = query.Where("user_uuid = ?", common.ClearQuery(userUuid))
	}

	if err := query.Find(&esims).Error; err != nil {
		return nil, err
	}

	var responseDTOs []dto.GetEsimDbDTO
	for _, esim := range esims {
		responseDTO := dto.GetEsimDbDTO{
			ID:            esim.ID,
			Iccid:         esim.Iccid,
			Lpa:           esim.Lpa,
			UserUUID:      esim.UserUUID,
			ProfileStatus: esim.ProfileStatus,
		}
		responseDTOs = append(responseDTOs, responseDTO)
	}

	return responseDTOs, nil
}
