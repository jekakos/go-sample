package packages

import (
	"sample-service/app/packages/dto"
	"sample-service/app/packages/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// @Summary Update Package
// @Description Update Package
// @Tags root
// @Accept application/json
// @Produce json
// @Param   package_id path string true "Package ID"
// @Param   request body dto.UpdatePackageRequestDTO true "Update package request payload"
// @Success 200     {object} dto.UpdatePackageResponseDTO
// @Router /package/{package_id} [patch]
func UpdatePackage(c *fiber.Ctx) error {

	println("Update Package")

	var data dto.UpdatePackageRequestDTO
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid data parameters")
	}

	packageId := c.Params("package_id")

	updatedPackage, err := repository.UpdatePackage(packageId, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update package")
	}

	region, err := repository.GetRegionAndCountries(updatedPackage.CountryList)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to show package")
	}

	dataGb := decimal.NewFromInt(updatedPackage.DataAmount).Div(decimal.NewFromInt(1000))

	dto := dto.UpdatePackageResponseDTO{
		Id:                  updatedPackage.GlobalId,
		DataGb:              dataGb,
		Days:                int64(updatedPackage.Duration),
		VendorCustomerPrice: updatedPackage.VendorCustomerPrice,
		CustomPrice:         updatedPackage.CustomPrice,
		VendorPrice:         updatedPackage.VendorPrice,
		Hidden:              updatedPackage.Hidden,
		RegionName:          region.Name,
		RegionCodeName:      region.CodeName,
	}

	return c.JSON(dto)

}
