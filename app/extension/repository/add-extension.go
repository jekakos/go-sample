package repository

import (
	"fmt"
	"sample-service/app/common"
	esimDto "sample-service/app/esim/dto"
	extDto "sample-service/app/extension/dto"
	"sample-service/app/infra"
	"sample-service/app/model"
)

func AddExtension(params extDto.AddExtensionDbDTO) (model.EsimAssignedPlan, error) {

	params.Iccid = common.ClearQuery(params.Iccid)

	var esim model.Esim
	if err := infra.DB.Where("iccid = ?", params.Iccid).First(&esim).Error; err != nil {
		return model.EsimAssignedPlan{}, fmt.Errorf("failed to find Esim with ICCID %s: %w", params.Iccid, err)
	}

	// Create the new EsimAssignedPlan
	newPlan := model.EsimAssignedPlan{
		EsimID:                   params.EsimID,
		PackageId:                params.PackageId,
		InitialQuantityInBytes:   params.InitialQuantityInBytes,
		RemainingQuantityInBytes: params.RemainingQuantityInBytes,
		StartTime:                params.StartTime,
		EndTime:                  params.EndTime,
		IsExpired:                params.IsExpired,
		VendorId:                 params.VendorId,
		VendorPlanId:             params.VendorPlanId,
		VendorAreaId:             params.VendorAreaId,
	}

	for _, areaDTO := range params.Areas {
		area := transformAreaDTOToArea(areaDTO)
		if err := infra.DB.FirstOrCreate(&area, model.Area{ID: areaDTO.Id}).Error; err != nil {
			return model.EsimAssignedPlan{}, fmt.Errorf("failed to find or create Area: %w", err)
		}
		newPlan.Areas = append(newPlan.Areas, area)
	}

	// Save the new plan
	if err := infra.DB.Create(&newPlan).Error; err != nil {
		return model.EsimAssignedPlan{}, fmt.Errorf("failed to create new EsimAssignedPlan: %w", err)
	}

	// Update the Esim's assigned plans if necessary
	esim.AssignedPlans = append(esim.AssignedPlans, newPlan)
	if err := infra.DB.Save(&esim).Error; err != nil {
		return model.EsimAssignedPlan{}, fmt.Errorf("failed to update Esim with new plan: %w", err)
	}

	return newPlan, nil

}

func transformAreaDTOToArea(areaDTO esimDto.AreaDTO) model.Area {
	return model.Area{
		ID:     areaDTO.Id,
		Name:   areaDTO.Name,
		Region: areaDTO.Region,
		Iso:    &areaDTO.Iso,
	}
}
