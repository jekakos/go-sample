package repository

import (
	"sample-service/app/common"
	"sample-service/app/extension/dto"
	"sample-service/app/infra"
)

func GetExtension(params dto.GetExtensionRequestDTO) ([]dto.GetExtensionDbDTO, error) {

	var extensions []dto.GetExtensionDbDTO
	query := infra.DB.Table(`esim_assigned_plans`).
		Select(`
			esim_assigned_plans.id, 
			esim_assigned_plans.esim_id, 
			esim_assigned_plans.package_id,
			esim_assigned_plans.initial_quantity_in_bytes, 
			esim_assigned_plans.remaining_quantity_in_bytes, 
			esim_assigned_plans.start_time, 
			esim_assigned_plans.end_time, 
			esim_assigned_plans.is_expired,
			esims.user_uuid,
			esims.iccid,
			esims.lpa`).
		Joins("INNER JOIN esims ON esim_assigned_plans.esim_id = esims.id")

	//Select("id, esim_id, initialQuantityInBytes, remainingQuantityInBytes, startTime, EndTime, IsExpired")
	iccid := common.GetStr(params.Iccid)
	id := common.GetInt32(params.ID)
	userUuid := common.GetStr(params.UserUUID)

	if id != 0 {
		query = query.Where("esim_assigned_plans.id = ?", id)
	}

	if iccid != "" {
		query = query.Where("esims.iccid = ?", common.ClearQuery(iccid))
	}

	if userUuid != "" {
		query = query.Where("esims.user_uuid = ?", common.ClearQuery(userUuid))
	}

	if err := query.Find(&extensions).Error; err != nil {
		return nil, err
	}

	return extensions, nil
}
