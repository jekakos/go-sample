package packages

import (
	"fmt"
	"sample-service/app/infra"
	vendor "sample-service/app/internal_vendor"
	vendorDto "sample-service/app/internal_vendor/dto"
	"sample-service/app/model"
	"sample-service/app/packages/dto"
	"sample-service/app/packages/repository"
)

func LoadPackagesService(reset string) (string, error) {

	packages, err := vendor.LoadPackages()
	if err != nil {
		println("Error while getting LoadPackages")
		return "", fmt.Errorf(err.Error())
	}

	vendorPackagesCount := len(packages)

	if reset == "1" || reset == "true" {
		// Если параметр reset присутствует, очищаем все три таблицы
		if err := infra.DB.Exec("TRUNCATE TABLE packages RESTART IDENTITY CASCADE").Error; err != nil {
			return "", fmt.Errorf("Failed to clear packages table")
		}
		if err := infra.DB.Exec("TRUNCATE TABLE areas RESTART IDENTITY CASCADE").Error; err != nil {
			return "", fmt.Errorf("Failed to clear areas table")
		}
		if err := infra.DB.Exec("TRUNCATE TABLE package_areas RESTART IDENTITY CASCADE").Error; err != nil {
			return "", fmt.Errorf("Failed to clear package_areas table")
		}
	}

	packagesAdded := 0
	packagesUpdated := 0
	packagesAreasUpdated := 0
	// This loop updates or adds
	for _, p := range packages {
		fmt.Printf("Package %d\n", p.Id)

		var existingPackage model.Package
		result := infra.DB.Preload("Areas").Where("vendor_id = ?", p.Id).Find(&existingPackage)

		if result.RowsAffected == 0 {

			areas := findOrCreateAreas(p.Areas)
			dataAmountInGB := p.DataAmount / 1000

			// Create Packeges
			packageData := model.Package{
				VendorId:    p.Id,
				GlobalId:    fmt.Sprintf("vendor-%d-%ddays-%dgb", p.Id, p.Duration, dataAmountInGB),
				DataAmount:  p.DataAmount,
				Duration:    p.Duration,
				VendorPrice: p.Price,
				Areas:       areas,
			}

			// Package was not found, create new
			if err := infra.DB.Create(&packageData).Error; err != nil {
				fmt.Printf("Error while create package %v\n", p)
				continue
			}
			packagesAdded++
		} else {
			// Found, Update
			countDbAreas := len(existingPackage.Areas)
			countVendorAreas := len(p.Areas)

			if !existingPackage.VendorPrice.Equal(p.Price) || countDbAreas != countVendorAreas {

				// Updae Packege

				if !existingPackage.VendorPrice.Equal(p.Price) {
					fmt.Printf("Vendor PRICE was changed %s -> %s\n", existingPackage.VendorPrice, p.Price)
					packageData := model.Package{
						VendorPrice: p.Price,
					}

					if err := infra.DB.Model(&existingPackage).Updates(packageData).Error; err != nil {
						fmt.Printf("Error while update package %v\n", p)
						continue
					}
				}

				// Update links

				if countDbAreas != countVendorAreas {
					fmt.Printf("Areas was changed %d -> %d\n", countDbAreas, countVendorAreas)
					areas := findOrCreateAreas(p.Areas)
					if err := infra.DB.Model(&existingPackage).Association("Areas").Replace(areas); err != nil {
						fmt.Printf("Error while update links %v\n", p)
						continue
					}
					packagesAreasUpdated++
				}
				packagesUpdated++
			}
		}
	}

	// This part marks as archive that packages that already not exists in Vendor response
	packagesDeleted := 0
	params := dto.GetPackageRequestDTO{}
	allDbPackages, err := repository.GetPackagesWithAreaCount(params)
	fmt.Printf("Detect deleted packages.\n Found %d db packages\n", len(allDbPackages))

	if err != nil {
		return "", fmt.Errorf("Failed to get packages after update")
	}

	for _, dbPackage := range allDbPackages {

		db_vendor_id := dbPackage.VendorId
		found := false
		for _, p := range packages {
			if p.Id == db_vendor_id {
				found = true
				continue
			}
		}

		if found {
			continue
		}

		fmt.Printf("Package %d was not found in Vendor response\n", db_vendor_id)
		// No db package in Vendor response -> Deleted = true
		var existingPackage model.Package
		infra.DB.Where("vendor_id = ?", db_vendor_id).Find(&existingPackage)
		// Update Packeges
		packageData := model.Package{
			Deleted: true,
		}
		infra.DB.Model(&existingPackage).Updates(packageData)
		packagesDeleted++
	}

	result := fmt.Sprintf("Vendor packages found: %d\nDB packages found: %d\nPackages added: %d\nPackages updated: %d\nPackages deleted: %d\nPackages Areas updated: %d\n",
		vendorPackagesCount,
		len(allDbPackages),
		packagesAdded,
		packagesUpdated,
		packagesDeleted,
		packagesAreasUpdated,
	)
	fmt.Printf(result)
	return result, nil

}

func findOrCreateAreas(vendorAreas []vendorDto.VendorArea) []model.Area {
	var areas []model.Area
	for _, a := range vendorAreas {

		var area model.Area
		if err := infra.DB.Where(&model.Area{Iso: &a.Iso}).FirstOrCreate(&area, model.Area{
			ID:     a.Id,
			Name:   a.Name,
			Region: a.Region,
			Iso:    &a.Iso,
		}).Error; err != nil {
			fmt.Printf("Error while find-or-create area %v\n", a)
			continue
		}
		areas = append(areas, area)
	}

	return areas
}
