package packages

import (
	"fmt"
	"sample-service/app/common"
	"sample-service/app/packages/dto"
	"sample-service/app/packages/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// @Summary Get Packages
// @Description Get Packages
// @Tags root
// @Accept application/json
// @Produce json
// @Param   request body dto.GetPackageRequestDTO true "Get package request payload"
// @Success 200     {object} []dto.GetPackageResponseDTO
// @Router /package [get]
func GetPackage(c *fiber.Ctx) error {

	var query dto.GetPackageRequestDTO
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid query parameters")
	}

	packages, err := repository.GetPackagesWithAreaCount(query)
	if err != nil {
		return common.ServerError(c, fmt.Sprintf("Error while getting GetPackagesWithAreaCount: %s", err.Error()))
	}

	var response []dto.GetPackageResponseDTO

	for _, p := range packages {

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

	return c.JSON(response)

}
