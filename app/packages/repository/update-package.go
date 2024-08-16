package repository

import (
	"fmt"
	"sample-service/app/common"
	"sample-service/app/infra"
	"sample-service/app/model"
	"sample-service/app/packages/dto"
)

func UpdatePackage(global_id string, data dto.UpdatePackageRequestDTO) (PackageWithAreaCount, error) {

	fmt.Printf("Data: %v", data)

	packageId, err := common.GetPackageFromGlobalId(global_id)
	fmt.Printf("PARAMS: %d", packageId)

	if err != nil {
		return PackageWithAreaCount{}, fmt.Errorf("Package ID parse error")
	}

	// Find package
	query := dto.GetPackageRequestDTO{
		Id: &global_id,
	}

	packages, err := GetPackagesWithAreaCount(query)

	if err != nil || len(packages) != 1 {
		return PackageWithAreaCount{}, fmt.Errorf("Cannot find package %s to update", global_id)
	}

	findPackage := packages[0]

	// Prepare map for updates
	updates := make(map[string]interface{})

	if data.VendorPrice != nil {
		updates["vendor_price"] = *data.VendorPrice
		findPackage.VendorPrice = *data.VendorPrice
	}

	if data.Days != nil {
		updates["duration"] = *data.Days
		findPackage.Duration = int16(*data.Days)
	}

	if data.Hidden != nil {
		updates["hidden"] = *data.Hidden
		findPackage.Hidden = *data.Hidden
	}

	if data.CustomPrice != nil {
		updates["custom_price"] = *data.CustomPrice
		findPackage.CustomPrice = data.CustomPrice
	} else if data.NullCustomPrice != nil && *data.NullCustomPrice {
		updates["custom_price"] = nil
		findPackage.CustomPrice = nil
	}

	// Update package in the database
	if err := infra.DB.Model(&model.Package{}).Where("vendor_id = ?", packageId).Updates(updates).Error; err != nil {
		return PackageWithAreaCount{}, err
	}

	return findPackage, nil
}
