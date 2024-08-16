package repository

import (
	"sample-service/app/common"
	"sample-service/app/esim/dto"
	"sample-service/app/infra"
	"sample-service/app/model"
)

func AddEsim(userUuid string, planId int64, packageId string, response dto.VendorCreateEsimResponseDTO) (model.Esim, error) {

	esim := model.Esim{
		Iccid:         response.Iccid,
		Lpa:           common.ExtractLpa(response.IosSetup),
		SmdpAddress:   response.SmdpAddress,
		ProfileStatus: response.ProfileStatus,
		UserUUID:      userUuid,
		InstalledAt:   common.ParseTime(response.InstalledAt),
	}

	// Create AssignedPlans from AssignedPlansDTO
	for _, planDTO := range response.AssignedPlans {
		plan := model.EsimAssignedPlan{
			VendorId:                 planDTO.Id,
			VendorPlanId:             planId,
			InitialQuantityInBytes:   planDTO.InitialQuantityInBytes,
			RemainingQuantityInBytes: planDTO.RemainingQuantityInBytes,
			IsExpired:                planDTO.IsExpired,
			StartTime:                common.ParseTime(planDTO.StartTime),
			EndTime:                  common.ParseTime(planDTO.EndTime),
			PackageId:                packageId,
		}

		// New EsimArea from AreaDTO
		for _, areaDTO := range planDTO.Areas {
			area := model.Area{
				ID:     areaDTO.Id,
				Name:   areaDTO.Name,
				Region: areaDTO.Region,
				Iso:    &areaDTO.Iso,
			}

			// Check if area exists, if not create it
			if err := infra.DB.Where(&model.Area{ID: area.ID}).FirstOrCreate(&area).Error; err != nil {
				return model.Esim{}, err
			}

			plan.Areas = append(plan.Areas, area)
		}

		esim.AssignedPlans = append(esim.AssignedPlans, plan)
	}

	// Сохранение объекта Esim в базе данных
	if err := infra.DB.Create(&esim).Error; err != nil {
		return model.Esim{}, err
	}

	// Link existing areas to the plans
	for _, plan := range esim.AssignedPlans {
		for _, area := range plan.Areas {
			if err := infra.DB.Model(&plan).Association("Areas").Append(&area); err != nil {
				return model.Esim{}, err
			}
		}
	}

	return esim, nil
}
