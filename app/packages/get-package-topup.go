package packages

import (
	"fmt"
	"sample-service/app/common"
	vendor "sample-service/app/internal_vendor"
	"sample-service/app/packages/dto"
	"sample-service/app/packages/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// @Summary Get Packages for Top-up
// @Description Get Packages for Top-up
// @Tags root
// @Accept application/json
// @Produce json
// @Param   iccid path string true "Iccid"
// @Success 200     {object} dto.GetPackageResponseDTO
// @Router /package/{iccid} [get]
func GetPackagesTopup(c *fiber.Ctx) error {

	iccid := c.Params("iccid")

	if !common.IsValidICCID(iccid) {
		return common.ServerError(c, fmt.Sprintf("Invalid ICCID format: %s", iccid))
	}

	// Get list of Packages from Vendor
	vendorPackages, err := vendor.GetPackagesForTopup(iccid)
	if err != nil {
		return common.ServerError(c, fmt.Sprintf("Error while getting GetPackagesTopup: %s", err.Error()))
	}

	// Generate IDS for db select
	var ids []int32

	for _, vendorPackage := range vendorPackages {
		ids = append(ids, int32(vendorPackage.Id))
	}

	// Select packages
	params := dto.GetPackageRequestDTO{
		VendorIds: &ids,
	}

	dbPackages, err := repository.GetPackagesWithAreaCount(params)

	// Generate response
	var response []dto.GetPackageResponseDTO

	for _, p := range dbPackages {

		region, err := repository.GetRegionAndCountries(p.CountryList)

		if err != nil {
			continue
		}

		dataGb := decimal.NewFromInt(p.DataAmount).Div(decimal.NewFromInt(1000))

		dto := dto.GetPackageResponseDTO{
			Id:                  p.GlobalId,
			DataGb:              dataGb,
			Days:                int64(p.Duration),
			VendorCustomerPrice: p.VendorCustomerPrice,
			CustomPrice:         p.CustomPrice,
			VendorPrice:         p.VendorPrice,
			Hidden:              p.Hidden,
			Region:              region,
		}

		response = append(response, dto)
	}

	c.Set("Content-Type", "application/json")
	return c.JSON(response)
}
