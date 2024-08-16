package packages

import (
	"sample-service/app/internal_vendor"

	"github.com/gofiber/fiber/v2"
)

func VendorPackagesHandler(c *fiber.Ctx) error {

	packages, err := internal_vendor.LoadPackages()
	if err != nil {
		println("Error while getting LoadPackages")
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	c.Set("Content-Type", "application/json")
	return c.JSON(packages)
}
